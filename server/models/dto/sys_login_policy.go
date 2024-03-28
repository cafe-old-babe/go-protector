package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/core/base"
	"go-protector/server/core/consts"
	"go-protector/server/core/current"
	"go-protector/server/core/custom/c_type"
)

// LoginPolicyProcessDTO 处理登录
type LoginPolicyProcessDTO struct {
	Service   base.Service
	PolicyDTO ILoginPolicyDTO
	ParamDTO  *LoginPolicyParamDTO
	User      *current.User
}

// LoginPolicyParamDTO 登录使用的参数
type LoginPolicyParamDTO struct {
	SessionId  string                 `json:"sessionId,omitempty"`
	PolicyCode c_type.LoginPolicyCode `json:"policyCode,omitempty"`
	Operate    int                    `json:"operate,omitempty"` // 0:发送邮箱 1:校验操作
	Val        string                 `json:"val,omitempty"`     // 校验值
}

// LoginPolicyResultDTO 登录使用的结果
type LoginPolicyResultDTO struct {
	SessionId  string                   `json:"sessionId"`
	Mode       string                   `json:"mode"` // 模式,0:通过全部 or 通过一个即可
	PolicyCode []c_type.LoginPolicyCode `json:"policyCode"`
}

var LoginPolicyDTOFactory = make(map[c_type.LoginPolicyCode]ILoginPolicyDTO)

func init() {
	LoginPolicyDTOFactory[consts.LoginPolicyGlobal] = &GlobalLoginPolicyDTO{}
	LoginPolicyDTOFactory[consts.LoginPolicyOtp] = &OTPLoginPolicyDTO{}
	LoginPolicyDTOFactory[consts.LoginPolicyEmail] = &EmailLoginPolicyDTO{}
	LoginPolicyDTOFactory[consts.LoginPolicyShare] = &ShareLoginPolicyDTO{}
	LoginPolicyDTOFactory[consts.LoginPolicyIntruder] = &IntruderLoginPolicyDTO{}
}

func NewLoginPolicyDTO(t c_type.LoginPolicyCode, param map[string]interface{}) (ILoginPolicyDTO, error) {
	dto, ok := LoginPolicyDTOFactory[t]
	if !ok {
		return nil, errors.New(fmt.Sprintf("无法识别: %v", t))
	}
	return dto.New(param)
}

// ILoginPolicyDTO 登录策略的接口
type ILoginPolicyDTO interface {
	New(param map[string]interface{}) (ILoginPolicyDTO, error)
	// IsEnable 是否启用
	IsEnable() bool
	// GetKey 通过Key获取value
	GetKey(key string) interface{}
	// Validate 校验
	Validate(dto ILoginPolicyDTO) error
}

// basePolicyDTO  ILoginPolicyDTO 的实现
type basePolicyDTO struct {
	rawParam map[string]interface{}
	// 私有的校验不生效
	Enable string `binding:"required,oneof=0 1"`
}

func (_self *basePolicyDTO) unmarshal(
	param map[string]interface{}, dto ILoginPolicyDTO) (ILoginPolicyDTO, error) {
	_self.rawParam = param
	_self.Enable = fmt.Sprintf("%v", param["enable"])
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, dto)
	return dto, err
}

func (_self *basePolicyDTO) Validate(dto ILoginPolicyDTO) error {
	if err := binding.Validator.ValidateStruct(_self); err != nil {
		return err
	}
	if err := binding.Validator.ValidateStruct(dto); err != nil {
		return err
	}
	return nil
}

func (_self *basePolicyDTO) IsEnable() bool {
	//enable, ok := _self.rawParam["enable"]
	//if !ok {
	//	return ok
	//}
	//var val int
	//valueOf := reflect.ValueOf(enable)
	//if valueOf.CanInt() {
	//	val = int(valueOf.Int())
	//} else {
	//	kind := valueOf.Kind()
	//	switch kind {
	//	case reflect.Bool:
	//		return valueOf.Bool()
	//	case reflect.String:
	//		val, _ = strconv.Atoi(valueOf.String())
	//	default:
	//		val = 0
	//	}
	//}
	return _self.Enable == "1"
}

func (_self *basePolicyDTO) GetKey(key string) interface{} {
	val, ok := _self.rawParam[key]
	if ok {
		return val
	}
	return nil
}

// OTPLoginPolicyDTO basePolicyDTO 的子类
// 如果一个子类对象可以用其父类对象替换，那么这个子类必须满足父类所有的接口。
// Go 语言的继承并不完全满足里氏替换原则
// https://github.com/Authenticator-Extension/Authenticator
type OTPLoginPolicyDTO struct {
	// Go 语言没有显式的继承语法。在 Go 语言中，要实现继承，
	// 需要使用嵌入式结构体（embedded struct）和指针。
	// 嵌入式结构体允许在一个结构体中嵌入另一个结构体的字段，从而实现继承的效果。
	basePolicyDTO `json:"-"`
	// Issuer 签发名称
	Issuer string `json:"issuer" binding:"required,max=64"`
	// Period 间隔时长
	Period uint `json:"period" binding:"required,min=30,max=60"`
	// SecretSize
	SecretSize uint `json:"secretSize" binding:"required,gte=12,max=24"`
}

func (_self OTPLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	return _self.basePolicyDTO.unmarshal(param, &_self)
	//_self.basePolicyDTO.rawParam = param
	//marshal, err := json.Marshal(param)
	//if err != nil {
	//	return nil, err
	//}
	//err = json.Unmarshal(marshal, &_self)
	//return _self, err
}

type EmailLoginPolicyDTO struct {
	basePolicyDTO `json:"-"`
	ExpireTime    uint `json:"expireTime" binding:"required,min=2,max=10"` // 分钟
}

func (_self EmailLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	return _self.basePolicyDTO.unmarshal(param, &_self)
	//var err error
	//_self.basePolicyDTO.rawParam = param
	//marshal, err := json.Marshal(param)
	//if err != nil {
	//	return nil, err
	//}
	//err = json.Unmarshal(marshal, &_self)
	//return _self, err
}

type GlobalLoginPolicyDTO struct {
	basePolicyDTO `json:"-"`
	Mode          string `json:"mode"  binding:"required,oneof=0 1"`
}

func (_self GlobalLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	return _self.basePolicyDTO.unmarshal(param, &_self)
	//var err error
	//_self.basePolicyDTO.rawParam = param
	//marshal, err := json.Marshal(param)
	//if err != nil {
	//	return nil, err
	//}
	//err = json.Unmarshal(marshal, &_self)
	//return _self, err
}

type ShareLoginPolicyDTO struct {
	basePolicyDTO `json:"-"`
	// 登录间隔
	LoginInterval uint `json:"loginInterval"  binding:"required;max=10,min=0"`
	// 单用户登陆 0:不限制;1;仅一人在线
	SingleOnline uint `json:"singleOnline"  binding:"required;oneof=0 1"`
	// 处理单用户登录 0:禁止当前用户登录,1:踢掉在线用户
	SingleOnlineOperate uint `json:"singleOnlineOperate"  binding:"required;oneof=0 1"`
}

func (_self ShareLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	return _self.basePolicyDTO.unmarshal(param, &_self)
	//var err error
	//_self.basePolicyDTO.rawParam = param
	//marshal, err := json.Marshal(param)
	//if err != nil {
	//	return nil, err
	//}
	//err = json.Unmarshal(marshal, &_self)
	//return _self, err
}

type IntruderLoginPolicyDTO struct {
	basePolicyDTO `json:"-"`
	FailThreshold uint `json:"failThreshold"  binding:"max=10,min=5"`
}

func (_self IntruderLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	return _self.basePolicyDTO.unmarshal(param, &_self)
}
