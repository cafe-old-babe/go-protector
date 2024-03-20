package service

import (
	"encoding/json"
	"fmt"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/consts/table_name"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/custom/c_type"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
)

type SysLoginPolicy struct {
	base.Service
}
type argType struct {
	loginPolicyType c_type.LoginPolicyCode
	assign          entity.SysLoginPolicy
}

var (
	defaultGlobalPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyGlobal,
		PolicyName: "全局登录策略",
		Enable:     "0",
	}
	defaultEmailPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyEmail,
		PolicyName: "邮箱校验策略",
		Enable:     "0",
	}
	defaultOtpPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyOtp,
		PolicyName: "动态密码校验策略",
		Enable:     "0",
	}
	defaultSharePolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyShare,
		PolicyName: "共享登录校验策略",
		Enable:     "0",
	}
	defaultIntruderPolicy = entity.SysLoginPolicy{
		PolicyCode: consts.LoginPolicyShare,
		PolicyName: "放爆破登录策略",
		Enable:     "0",
	}
	infoArgs []argType
)

func init() {
	defaultParam := map[string]interface{}{
		"mode": "1", // 0: 校验全部,1:只校验一个策略即可
	}
	bytes, _ := json.Marshal(defaultParam)
	defaultGlobalPolicy.Json = string(bytes)

	defaultParam = map[string]interface{}{
		"expireTime": "5", // 分钟
	}
	bytes, _ = json.Marshal(defaultParam)
	defaultEmailPolicy.Json = string(bytes)

	defaultParam = map[string]interface{}{
		"issuer":     "issuer", // 签发名称
		"period":     5,        // 秒
		"secretSize": 24,       // 秘钥长度
	}
	bytes, _ = json.Marshal(defaultParam)
	defaultOtpPolicy.Json = string(bytes)
	bytes, _ = json.Marshal(dto.ShareLoginPolicyDTO{LoginInterval: 0})
	defaultSharePolicy.Json = string(bytes)

	bytes, _ = json.Marshal(dto.IntruderLoginPolicyDTO{
		FailThreshold: 5,
		LockTime:      5,
	})
	defaultIntruderPolicy.Json = string(bytes)

	infoArgs = []argType{
		{
			loginPolicyType: consts.LoginPolicyGlobal,
			assign:          defaultGlobalPolicy,
		}, {
			loginPolicyType: consts.LoginPolicyEmail,
			assign:          defaultEmailPolicy,
		}, {
			loginPolicyType: consts.LoginPolicyOtp,
			assign:          defaultOtpPolicy,
		},
		{
			loginPolicyType: consts.LoginPolicyShare,
			assign:          defaultSharePolicy,
		}, {

			loginPolicyType: consts.LoginPolicyIntruder,
			assign:          defaultIntruderPolicy,
		},
	}
}

// Info 查询策略信息
func (_self *SysLoginPolicy) Info() (result *dto.Result) {
	info := map[c_type.LoginPolicyCode]map[string]interface{}{}
	var modelSlice []entity.SysLoginPolicy
	if err := _self.DB.Table(table_name.SysLoginPolicy).Find(&modelSlice).Error; err != nil {
		return dto.ResultFailureErr(err)
	}

	var model entity.SysLoginPolicy
	var valMap map[string]interface{}
	if len(modelSlice) == len(dto.LoginPolicyMap) {
		for _, model = range modelSlice {
			valMap = make(map[string]interface{})
			if err := json.Unmarshal([]byte(model.Json), &valMap); err != nil {
				return dto.ResultFailureErr(err)
			}
			valMap["name"] = model.PolicyName
			valMap["enable"] = model.Enable
			info[model.PolicyCode] = valMap
			model = entity.SysLoginPolicy{}
		}
	} else {
		for _, arg := range infoArgs {
			if err := _self.DB.Where(entity.SysLoginPolicy{PolicyCode: arg.loginPolicyType}).
				Attrs(arg.assign).FirstOrCreate(&model).Error; err != nil {
				return dto.ResultFailureErr(err)
			}
			valMap = make(map[string]interface{})
			if err := json.Unmarshal([]byte(model.Json), &valMap); err != nil {
				return dto.ResultFailureErr(err)
			}
			valMap["name"] = model.PolicyName
			valMap["enable"] = model.Enable
			info[arg.loginPolicyType] = valMap
			model = entity.SysLoginPolicy{}
		}
	}

	/*
		if err := _self.DB.Where(entity.SysLoginPolicy{PolicyCode: consts.LoginPolicyGlobal}).
			Assign(defaultGlobalPolicy).FirstOrCreate(&model).Error; err != nil {
			return dto.ResultFailureErr(err)
		}
		if err := json.Unmarshal([]byte(model.PJson), &valMap); err != nil {
			return dto.ResultFailureErr(err)
		}
		valMap["name"] = model.PolicyName
		info[consts.LoginPolicyGlobal] = valMap
		clear(valMap)
		model = entity.SysLoginPolicy{}

		if err := _self.DB.Where(entity.SysLoginPolicy{PolicyCode: consts.LoginPolicyEmail}).
			Assign(defaultEmailPolicy).FirstOrCreate(&model).Error; err != nil {
			return dto.ResultFailureErr(err)
		}
		if err := json.Unmarshal([]byte(model.PJson), &valMap); err != nil {
			return dto.ResultFailureErr(err)
		}
		valMap["name"] = model.PolicyName
		info[consts.LoginPolicyEmail] = valMap
		clear(valMap)
	*/
	return dto.ResultSuccess(info)
}

// Save 保存
func (_self *SysLoginPolicy) Save(param map[c_type.LoginPolicyCode]map[string]interface{}) (
	result *dto.Result) {

	if param == nil || len(param) <= 0 {
		return dto.ResultFailureErr(c_error.ErrParamInvalid)
	}
	// 只保存是否启用和json
	db := _self.DB.Begin()

	defer func() {
		if result.IsSuccess() {
			db.Commit()
			return
		}
		db.Rollback()
	}()

	for k, mapVal := range param {
		// 默认关闭
		enable := "0"
		if val, ok := mapVal["enable"]; ok {
			//interface 转 int
			enable = fmt.Sprintf("%v", val)
		}
		marshal, err := json.Marshal(mapVal)
		if err != nil {
			return dto.ResultFailureErr(err)
		}

		updateMap := map[string]interface{}{
			"enable": enable,
			"json":   string(marshal),
		}
		policyDTO, err := dto.NewLoginPolicyDTO(k, mapVal)
		if err = policyDTO.Validate(policyDTO); err != nil {
			return dto.ResultFailureErr(err)
		}
		if err = db.Table(table_name.SysLoginPolicy).
			Where("policy_code = ?", k).Updates(updateMap).Error; err != nil {
			return dto.ResultFailureErr(err)
		}
	}
	return dto.ResultSuccessMsg("保存成功")

}

// ProcessLoginPolicy 处理登录策略
func (_self *SysLoginPolicy) ProcessLoginPolicy(loginDTO dto.Login, user *entity.SysUser) (res *dto.Result) {

	return
}
