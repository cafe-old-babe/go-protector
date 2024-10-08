package service

import (
	"errors"
	"fmt"
	"go-protector/server/biz/model/dao"
	"go-protector/server/biz/model/dto"
	"go-protector/server/biz/model/entity"
	"go-protector/server/internal/base"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/consts/table_name"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/custom/c_type"
	"go-protector/server/internal/database/condition"
	"go-protector/server/internal/ssh/ssh_con"
	"go-protector/server/internal/utils/async"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"time"
)

var findAssetInfoAccountSliceMapByCollectors = map[string]func(base.IService, []uint64) ([]entity.AssetInfoAccount, error){}
var findAssetAccountInfoSliceMapByDial = map[string]func(base.IService, []uint64) ([]entity.AssetAccountInfo, error){}

func init() {
	findAssetInfoAccountSliceMapByCollectors["asset"] = func(_self base.IService, ids []uint64) (slice []entity.AssetInfoAccount, err error) {
		return findAssetInfoAccountSliceByIds(_self, ids)
	}
	findAssetInfoAccountSliceMapByCollectors["account"] = func(_self base.IService, ids []uint64) ([]entity.AssetInfoAccount, error) {
		var accountService AssetAccount
		_self.MakeService(&accountService)
		return accountService.FindAssetInfoAccountSliceByIds(ids)
	}

	// 拨测资产特权账号
	findAssetAccountInfoSliceMapByDial["asset"] = func(_self base.IService, ids []uint64) (slice []entity.AssetAccountInfo, err error) {

		return findAssetAccountInfoSliceByIds(_self, ids)
	}

	// 拨测指定账号
	findAssetAccountInfoSliceMapByDial["account"] = func(_self base.IService, ids []uint64) (slice []entity.AssetAccountInfo, err error) {
		var accountService AssetAccount
		_self.MakeService(&accountService)
		return accountService.FindAssetAccountInfoSliceByIds(ids)
	}

}

type AssetInfo struct {
	base.Service
}

func (_self *AssetInfo) Page(req *dto.AssetInfoPageReq) (res *base.Result) {
	var err error
	var assetInfoSlice []*entity.AssetInfo
	var count int64
	defer func() {
		if res != nil {
			return
		}
		if err != nil {
			res = base.ResultFailureErr(err)
			return
		}
		res = base.ResultPage(assetInfoSlice, req, count)
		return
	}()
	if req == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	tx := _self.GetDB().Scopes(
		condition.Paginate(req),
		condition.Like(table_name.AssetBasic+".asset_name", req.AssetName),
		condition.Like(table_name.AssetBasic+".IP", req.IP),
		func(db *gorm.DB) *gorm.DB {
			if len(req.GroupIds) > 0 {
				db = db.Where(table_name.AssetBasic+".asset_group_id in ?", req.GroupIds)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB {
			if !req.Auth {
				return db.Joins("RootAcc", _self.GetDB().Where("account_type = ?", "0"))
			}
			return db.Where(table_name.AssetBasic+".id in (?)",
				dao.AssetAuth.GenerateSubQueryForAssetId(db, current.GetUserId(_self.GetContext())))

		},
	)

	err = tx.Joins("AssetGateway").
		Joins("ManagerUser").
		Joins("AssetGroup").
		//Joins("left join " + table_name.AssetAccount +
		//	" on " + table_name.AssetBasic + ".id = " + table_name.AssetAccount + ".asset_id and account_type = '0' ").
		Find(&assetInfoSlice).
		Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		return
	}
	var accountService AssetAccount
	_self.MakeService(&accountService)

	// 清空密码
	for i := range assetInfoSlice {
		assetInfoSlice[i].RootAcc.Password = ""
	}
	return
}

// Save 保存资产信息
func (_self *AssetInfo) Save(req *dto.AssetInfoSaveReq) (res *base.Result) {
	_self.Begin()
	defer func() {
		if err := recover(); err != nil {
			res = base.ResultFailureMsg(fmt.Sprintf("%s", err))
			goto rollback
		}
		if res.IsSuccess() {
			_self.GetDB().Commit()
			return
		}
	rollback:
		_self.GetDB().Rollback()
	}()
	// 校验资产信息
	assetBasic := entity.AssetBasic{
		ModelId:        entity.ModelId{ID: req.ID},
		AssetName:      req.AssetName,
		AssetGroupId:   req.GroupId,
		IP:             req.IP,
		Port:           req.Port,
		AssetGatewayId: req.AssetGatewayId,
		ManagerUserId:  req.ManagerUserId,
	}
	//_self.GetContext().Request = _self.GetContext().Request.WithContext(context.WithValue(_self.GetContext().Request.ctx(), "mutex", assetBasic))
	//_self.GetDB().WithContext(context.WithValue(_self.GetDB().Statement.ctx, "mutex", assetBasic))
	res = _self.SimpleSave(&assetBasic, func() error {
		return dao.AssetBasic.CheckSave(_self.GetDB(), assetBasic)

	})
	if !res.IsSuccess() {
		return
	}
	//privilegeAccount := req.PrivilegeAccount
	accountModel := entity.AssetAccount{
		Account:       req.Account,
		Password:      req.Password,
		AccountType:   "0",
		AccountStatus: "0",
		AssetId:       assetBasic.ID,
	}
	var accountService AssetAccount
	_self.MakeService(&accountService)
	res = accountService.Save(&accountModel)

	return
}

func (_self *AssetInfo) Delete(idsReq *base.IdsReq) (res *base.Result) {
	if idsReq == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	var ids []uint64
	if ids = idsReq.GetIds(); len(ids) <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	if err := _self.GetDB().Transaction(func(tx *gorm.DB) (err error) {
		// 删除资产
		if err = tx.Delete(&entity.AssetBasic{}, ids).Error; err != nil {
			return
		}
		// 删除从账号
		if err = dao.AssetAccount.DeleteByAssetId(tx, ids); err != nil {
			return
		}
		// 删除授权
		var auth entity.AssetAuth
		err = auth.DeleteRedundancy(tx, ids, entity.TypeAssetBasic)

		return
	}); err != nil {
		res = base.ResultFailureErr(err)

	}

	return
}

// Collectors 采集资产信息
func (_self *AssetInfo) Collectors(idsReq *base.IdsReq, collType string) (res *base.Result) {
	if len(collType) <= 0 || idsReq == nil {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}
	ids := idsReq.GetIds()
	if len(ids) <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
	}

	var collectorsFunc func(base.IService, []uint64) ([]entity.AssetInfoAccount, error)
	var ok bool
	if collectorsFunc, ok = findAssetInfoAccountSliceMapByCollectors[collType]; !ok {
		collectorsFunc = func(base.IService, []uint64) ([]entity.AssetInfoAccount, error) {
			return nil, c_error.ErrParamInvalid
		}
	}

	/* 1.0
	switch collType {
	case "asset":
		collectorsFunc = func() (slice []entity.AssetInfoAccount, err error) {
			err = _self.GetDB().Joins("Accounts").Find(&slice).Error
			return
		}
	case "account":
		collectorsFunc = func() ([]entity.AssetInfoAccount, error) {
			var accountService AssetAccount
			_self.MakeService(&accountService)
			return accountService.FindAssetInfoAccountSliceByIds(ids)
		}
	default:
		collectorsFunc = func() ([]entity.AssetInfoAccount, error) {
			return nil, c_error.ErrParamInvalid
		}
	}
	slice, err := collectorsFunc()
	*/

	slice, err := collectorsFunc(_self, ids)

	if err != nil {
		res = base.ResultFailureErr(err)
		return
	}

	res = _self.DoBatchCollectors(slice)

	// 添加采集任务

	return
}

func (_self *AssetInfo) FindAssetInfoAccountSliceByIds(ids []uint64) (slice []entity.AssetInfoAccount, err error) {
	return findAssetInfoAccountSliceByIds(_self, ids)
}

func findAssetInfoAccountSliceByIds(_self base.IService, ids []uint64) (slice []entity.AssetInfoAccount, err error) {
	if ids == nil || len(ids) <= 0 {
		return
	}
	err = _self.GetDB().Joins("AssetGateway").
		Joins("RootAcc", _self.GetDB().Where("RootAcc.account_type = ?", "0")).
		Preload("Accounts", _self.GetDB().Where("account_type <> ?", "0")).
		Find(&slice, ids).Error
	return
}

func findAssetAccountInfoSliceByIds(self base.IService, ids []uint64) (slice []entity.AssetAccountInfo, err error) {
	tx := self.GetDB()
	err = tx.Joins("AssetBasic").
		//Joins("left join " + table_name.AssetAccount +
		//	" on " + table_name.AssetBasic + ".id = " + table_name.AssetAccount + ".asset_id and account_type = '0' ").
		Find(&slice, "asset_account.account_type = '0' and AssetBasic.id in ? ", ids).Error

	return
}

// DoBatchCollectors 开始采集
func (_self *AssetInfo) DoBatchCollectors(slice []entity.AssetInfoAccount) (res *base.Result) {
	if slice == nil || len(slice) <= 0 {
		res = base.ResultFailureMsg("无采集信息")
		return
	}
	if slice == nil || len(slice) <= 0 {
		res = base.ResultFailureErr(c_error.ErrParamInvalid)
		return
	}

	//_self.GetDB() = _self.GetDB().WithContext(context.Background())
	_self.WithGoroutineDB()
	for _, assetInfo := range slice {
		//async.CommonWork.Submit(func() {
		//	_self.DoCollectors(assetInfo)
		//})
		async.CommonWorkPool.Submit(func() {
			_self.DoCollectors(assetInfo)
		})

	}
	return base.ResultSuccessMsg("采集任务下发完成")

}

func (_self *AssetInfo) Dial(idsReq *base.IdsReq, dialType string) (res *base.Result) {
	var f func(base.IService, []uint64) ([]entity.AssetAccountInfo, error)
	ids := idsReq.GetIds()
	var ok bool
	if f, ok = findAssetAccountInfoSliceMapByDial[dialType]; !ok {
		f = func(base.IService, []uint64) ([]entity.AssetAccountInfo, error) {
			return nil, c_error.ErrParamInvalid
		}
	}
	var slice []entity.AssetAccountInfo
	var err error
	if slice, err = f(_self, ids); err != nil {
		res = base.ResultFailureErr(err)
		return
	}

	return _self.DoBatchDail(slice)
}

func (_self *AssetInfo) DoBatchDail(slice []entity.AssetAccountInfo) (res *base.Result) {

	if len(slice) <= 0 {
		res = base.ResultFailureErr(errors.New("无拨测信息"))
		return
	}
	_self.WithGoroutineDB()
	for _, elem := range slice {
		async.CommonWorkPool.Submit(func() {
			_self.DoDail(elem)
		})
	}

	return base.ResultSuccessMsg("拨测任务下发完成")

}

func (_self *AssetInfo) DoDail(elem entity.AssetAccountInfo) {

	var client *ssh_con.Client

	var err error

	defer func() {
		if client != nil {
			if closeErr := client.Close(); closeErr != nil {
				_self.GetLogger().Error("close err: %v", closeErr)
			}
		}
		dailStatus := "1"
		nowTime := c_type.NowTime()
		dailMsg := "[" + nowTime.String() + "]"
		if err != nil {
			dailStatus = "0"
			dailMsg += err.Error()
		} else {
			dailMsg += "拨测成功"
		}

		if err = _self.GetDB().Updates(&entity.AssetAccount{
			ModelId:    elem.ModelId,
			DailStatus: dailStatus,
			DailMsg:    dailMsg,
		}).Error; err != nil {
			_self.GetLogger().Error("update dial result failure: %v", err)
		}
	}()

	client, err = ssh_con.Connect(&ssh_con.ConnectParam{
		IP:        elem.AssetBasic.IP,
		GatewayId: elem.AssetBasic.AssetGatewayId,
		Port:      elem.AssetBasic.Port,
		Username:  elem.Account,
		Password:  elem.Password,
		Timeout:   time.Second * 3,
	})

	return
}

// DoCollectors 采集
func (_self *AssetInfo) DoCollectors(assetInfo entity.AssetInfoAccount) {

	var err error
	var cli *ssh_con.Client
	defer func() {
		if cli != nil {
			_ = cli.Close()
		}

	}()
	if len(assetInfo.Accounts) <= 0 {
		return
	}
	//if deStr, err = gm.Sm4DecryptCBC(assetInfo.RootAcc.Password); err != nil {
	//	accountCollectorsDTO = append(accountCollectorsDTO, dto.AccountAnalysisExtendDTO{
	//		AssetId: assetInfo.ID,
	//		Err:     errors.Join(errors.New("密码解密失败"), err),
	//	})
	//	continue
	//}

	var accountCollectorsDTO []dto.AccountAnalysisExtendDTO
	if cli, err = ssh_con.Connect(&ssh_con.ConnectParam{
		ID:        assetInfo.ID,
		GatewayId: assetInfo.AssetGatewayId,
		IP:        assetInfo.IP,
		Port:      assetInfo.Port,
		Username:  assetInfo.RootAcc.Account,
		Password:  assetInfo.RootAcc.Password,
		Timeout:   0,
	}); err != nil {
		accountCollectorsDTO = append(accountCollectorsDTO, dto.AccountAnalysisExtendDTO{
			AssetId: assetInfo.ID,
			Err:     errors.Join(errors.New("连接失败"), err),
		})
		goto saveLabel
	}

	if err != nil {
		// 保存错误信息
		accountCollectorsDTO = append(accountCollectorsDTO, dto.AccountAnalysisExtendDTO{
			AssetId: assetInfo.ID,
			Err:     errors.Join(errors.New("创建会话失败"), err),
		})
		goto saveLabel
	}

	for _, account := range assetInfo.Accounts {
		// -1收集从账号,0-特权账号 不采集
		if account.AccountType == "-1" {
			continue
		}
		collectorsDTO := collectorsAccount(account, cli.SSHClient)
		collectorsDTO.AssetId = assetInfo.ID
		accountCollectorsDTO = append(accountCollectorsDTO, collectorsDTO)
	}
saveLabel:
	var accountService AssetAccount
	_self.MakeService(&accountService)

	accountService.AnalysisExtend(accountCollectorsDTO)
}

// collectorsAccount 采集从帐号
func collectorsAccount(account entity.AssetAccount, cli *ssh.Client) (collectorsDTO dto.AccountAnalysisExtendDTO) {
	collectorsDTO.ID = account.ID
	collectorsDTO.In = fmt.Sprintf(consts.CollFmt, account.Account, account.Account)

	var session *ssh.Session
	var err error
	if session, err = cli.NewSession(); err != nil {
		collectorsDTO.Err = err
		return
	}

	if collectorsDTO.Out, collectorsDTO.Err = session.Output(collectorsDTO.In); collectorsDTO.Err != nil {
		collectorsDTO.Err = errors.Join(errors.New("执行命令失败"), collectorsDTO.Err)
	}
	return
}
