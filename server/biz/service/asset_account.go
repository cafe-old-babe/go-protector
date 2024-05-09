package service

import (
	"errors"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type AssetAccount struct {
	base.Service
}

// Page 分页查询
func (_self *AssetAccount) Page(req *dto.AssetAccountPageReq) (res *base.Result) {
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var slice []entity.AssetAccount
	tx := _self.DB.Omit(table_name.AssetAccount+".password").Scopes(
		condition.Paginate(req),
		condition.Like(table_name.AssetAccount+".account", req.Account),
		condition.Like("AssetBasic.ip", req.IP),
		condition.Like("AssetBasic.asset_name", req.AssetName),
	)
	count, err := _self.Count(
		tx.Joins("AssetBasic").
			Joins("Extend").
			Order("AssetBasic.created_at desc").
			Order(table_name.AssetAccount + ".account_status").
			Order(table_name.AssetAccount + ".account_type").
			Find(&slice))

	if err != nil {
		res = base.ResultFailureErr(err)
	} else {
		res = base.ResultPage(slice, req, count)
	}
	return
}

func (_self *AssetAccount) FindAssetInfoAccountSliceByIds(ids []uint64) (slice []entity.AssetInfoAccount, err error) {
	if ids == nil || len(ids) <= 0 {
		return
	}
	var accountSlice []entity.AssetAccount
	if err = _self.DB.Find(&accountSlice, ids).Error; err != nil || len(accountSlice) <= 0 {
		return
	}
	var assetIdSlice []uint64
	var ok bool
	accountMapByAssetId := make(map[uint64][]entity.AssetAccount)
	for _, account := range accountSlice {
		if _, ok = accountMapByAssetId[account.AssetId]; !ok {
			assetIdSlice = append(assetIdSlice, account.AssetId)
		}
		accountMapByAssetId[account.AssetId] = append(accountMapByAssetId[account.AssetId], account)
	}

	if err = _self.DB.Joins("AssetGateway").Find(&slice, assetIdSlice).Error; err != nil || len(slice) <= 0 {
		return
	}

	for i := range slice {
		if accountSlice, ok = accountMapByAssetId[slice[i].ID]; ok {
			slice[i].Accounts = append(slice[i].Accounts[:0], accountSlice...)
		}
	}
	return
}

// FindRootAccMapByAssetIdSlice 根据资产ID 查询特权帐号
func (_self *AssetAccount) FindRootAccMapByAssetIdSlice(asseIdSlice []uint64) (resMap map[uint64]entity.AssetAccount, err error) {

	if asseIdSlice == nil || len(asseIdSlice) <= 0 {
		resMap = make(map[uint64]entity.AssetAccount)
		return
	}
	var accountSlice []entity.AssetAccount
	if err = _self.DB.Where("account_type = '0' and asset_id in ?", accountSlice).Find(&accountSlice).Error; err != nil {
		return
	}

	resMap = utils.SliceToUint64Map[entity.AssetAccount](accountSlice,
		func(elem entity.AssetAccount) uint64 {
			return elem.AssetId
		}, func(elem entity.AssetAccount) entity.AssetAccount {
			return elem
		})

	return
}

// FillAssetInfoRootAcc 填充特权帐号信息
// assetInfoSlice *[]entity.AssetInfo	elem 不能修改	(*assetInfoSlice)[i] 可以修改
// assetInfoSlice []entity.AssetInfo	elem 不能修改	assetInfoSlice[i] 可以修改
// assetInfoSlice []*entity.AssetInfo	elem 可以修改	assetInfoSlice[i] 可以修改
func (_self *AssetAccount) FillAssetInfoRootAcc(assetInfoSlice []*entity.AssetInfo) (err error) {
	if assetInfoSlice == nil || len(assetInfoSlice) <= 0 {
		return
	}
	var rootAccMap map[uint64]entity.AssetAccount
	asseIdSlice := utils.SliceToFieldSlice[uint64]("ID", assetInfoSlice)
	if rootAccMap, err = _self.FindRootAccMapByAssetIdSlice(asseIdSlice); err != nil {
		return
	}
	if len(rootAccMap) <= 0 {
		return
	}
	for _, elem := range assetInfoSlice {
		elem.RootAcc, _ = rootAccMap[elem.ID]
	}
	return
}

// AnalysisExtend 解析并保存从账号扩展信息
func (_self *AssetAccount) AnalysisExtend(dtoSlice []dto.AccountAnalysisExtendDTO) (res *base.Result) {

	var err error
	db := _self.DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
			if res == nil {
				res = base.ResultFailureErr(err)
			}
		} else {
			db.Commit()
			if res == nil {
				res = base.ResultSuccessMsg("解析成功")
			}
		}

	}()
	var account entity.AssetAccount
	var extend entity.AssetAccountExtend
	var outSlice []string
	var line []string
	var temp int
	zeroTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	nowTime := c_type.NowTime()

	//-1-采集失败,0-未采集信息,1-正常,2-即将过期,3-已过期,4-已禁用
	var accountStatus string
	for _, elem := range dtoSlice {
		if elem.Err != nil {
			if err = db.Model(&extend).Where("id in (?)",
				_self.DB.Model(&account).Select("id").
					Where("asset_id = ? and asset_type <> '0'"),
			).Updates(entity.AssetAccountExtend{
				CollectMsg:  elem.Err.Error(),
				CollectTime: nowTime,
			}).Error; err != nil {
				return
			}
			continue
		}

		if elem.ID <= 0 {
			continue
		}
		extend = entity.AssetAccountExtend{}
		extend.ID = elem.ID
		extend.CollectTime = nowTime
		outSlice = strings.Split(string(elem.Out), "\n")
		//redis:x:986:985:Redis Database Server:/var/lib/redis:/sbin/nologin
		//0-用户名（redis）：这是用户的登录名。
		//1-密码（x）：在/etc/shadow文件中存储了加密后的密码。
		//2-用户ID（986）：这是用户的唯一标识符。
		//3-组ID（985）：这是用户所属的主要组的唯一标识符。
		//4-用户描述（Redis Database Server）：这是用户的描述信息。
		//5-主目录（/var/lib/redis）：这是用户的主目录路径。
		//6-登录Shell（/sbin/nologin）：这是用户登录后默认使用的Shell。
		if line = strings.Split(outSlice[0], ":"); len(line) > 0 {
			extend.Uid = line[2]
			extend.HomePath = line[5]
			extend.Shell = line[6]
		}

		if len(outSlice) < 2 || len(outSlice[1]) <= 0 {
			continue
		}
		// username:password:lastchg:min:max:warn:inactive:expire:hold
		//1-username：账户名与/etc/passwd里面的账户名是一一对应的关系。
		//2-password：分为3类，分别是奇奇怪怪的字符串、*和!!其中，奇奇怪怪的字符串就是加密过的密码文件。星号代表帐号被锁定，双叹号表示这个密码已经过期了。
		//			 以$6$开头的，表明是用SHA-512加密的，$1$ 表明是用MD5加密的、$2$ 是用Blowfish加密的、$5$是用 SHA-256加密的。
		//3-lastchg：最后一次更改密码的日期：表示为自1970年1月1日00:00 UTC以来的天数。值0有一个特殊的含义，即用户下次登录系统时应该更改密码。空字段表示密码修改周期功能被禁用
		//4-min：这个是表明上一次修改密码的日期与1970-1-1相距的天数密码不可改的天数：假如这个数字是8，则8天内不可改密码，如果是0，则随时可以改。
		//5-max：如果是99999则永远不用改。如果是其其他数字比如12345，那么必须在距离1970-1-1的12345天内修改密码，否则密码失效。
		//6-warn：比如你在第五条规定今年6月20号规定密码必须被修改，系统会从距离6月20号的N天前向对应的用户发出警告。
		//7-inactive：假设这个数字被设定为M，那么帐号过期的M天内修改密码是可以修改的，改了之后账户可以继续使用。
		//8-expire：假设这个日期为X，与第三条一样，X表示的日期依然是1970-1-1相距的天数，过了X之后，帐号失效。
		//9-保留：被保留项，暂时还没有被用上。

		if line = strings.Split(outSlice[1], ":"); len(line) <= 0 {
			goto saveLabel
		}

		switch line[1] {
		case "*":
			accountStatus = "4"
		case "!", "!!":
			accountStatus = "3"
		default:
			accountStatus = "1"
		}

		// len(line[2]) > 0  新创建账号没有最后修改时间
		if temp, err = strconv.Atoi(line[2]); len(line[2]) > 0 && err == nil && temp > 0 {
			// 最后修改密码时间
			extend.LastUpPwdTime = c_type.NewTime(zeroTime.Add(time.Duration(temp) * 24 * time.Hour))
		}
		if !strings.HasSuffix(extend.Shell, "nologin") {
			goto saveLabel
		}

		if temp, err = strconv.Atoi(line[4]); err == nil || temp < 99999 {
			// 密码过期时间
			extend.ExpirePwdTime = c_type.NewTime(zeroTime.Add(time.Duration(temp) * 24 * time.Hour))
		} else if err == nil && temp >= 99999 {
			goto saveLabel
		}

		if temp, err = strconv.Atoi(line[6]); err == nil || extend.ExpirePwdTime.Valid {
			// 缓冲期至
			extend.InactiveAt = c_type.NewTime(extend.ExpirePwdTime.Time.Add(time.Duration(temp) * 24 * time.Hour))
		}
		if temp, err = strconv.Atoi(line[7]); err == nil {
			// 账号过期时间(锁定时间)
			extend.ExpireAt = c_type.NewTime(zeroTime.Add(time.Duration(temp) * 24 * time.Hour))

		}
		if temp, err = strconv.Atoi(line[5]); accountStatus == "1" && err == nil && extend.ExpirePwdTime.Valid {
			// 判断密码是否即将过期
			if time.Now().After(extend.ExpirePwdTime.Time.Add(-time.Duration(temp) * 24 * time.Hour)) {
				accountStatus = "2"
			}
		}

	saveLabel:
		err = db.Transaction(func(tx *gorm.DB) (err error) {
			// 保存从账号信息
			if err = tx.Where("id = ?", elem.ID).Updates(entity.AssetAccount{
				AccountStatus: accountStatus,
			}).Error; err != nil {
				return
			}
			err = tx.Where("id = ?", elem.ID).Updates(extend).Error
			return
		})

	}
	return

}

// Save 保存
func (_self *AssetAccount) Save(model *entity.AssetAccount) (res *base.Result) {
	res = _self.SimpleSave(model, func() error {
		return dao.AssetAccount.CheckSave(_self.DB, model)
	})
	return
}

// CheckBatchDeleteByIds 批量删除检查
func (_self *AssetAccount) CheckBatchDeleteByIds(ids []uint64) (err error) {
	if len(ids) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	var count int64
	count, err = _self.Count(_self.DB.Model(&entity.AssetAccount{}).Where("account_type = '0' and id in ?", ids))
	if err != nil {
		return
	}
	if count > 0 {
		err = errors.New("选中的从帐号中包括特权帐号,请核对后在操作")
	}
	return
}

// FindAssetAccountInfoSliceByIds 根据ID查询从帐号信息
func (_self *AssetAccount) FindAssetAccountInfoSliceByIds(ids []uint64) (slice []entity.AssetAccountInfo, err error) {
	if len(ids) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	err = _self.DB.Preload("AssetBasic").Find(&slice, ids).Error
	return
}
