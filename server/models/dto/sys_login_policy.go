package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"go-protector/server/core/consts"
	"go-protector/server/core/custom/c_type"
	"reflect"
	"strconv"
)

var codeMap = map[c_type.LoginPolicyCode]ILoginPolicyDTO{}

func init() {
	codeMap[consts.LoginPolicyGlobal] = &GlobalLoginPolicyDTO{}
	codeMap[consts.LoginPolicyOtp] = &OTPLoginPolicyDTO{}
	codeMap[consts.LoginPolicyEmail] = &EmailLoginPolicyDTO{}
}

func CreateLoginPolicyDTO(t c_type.LoginPolicyCode, param map[string]interface{}) (ILoginPolicyDTO, error) {
	dto, ok := codeMap[t]
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
	// Verify 校验
	Verify() error
}

// commonLoginPolicyDTO  ILoginPolicyDTO 的实现
type commonLoginPolicyDTO struct {
	rawParam map[string]interface{}
}

func (_self commonLoginPolicyDTO) IsEnable() bool {
	enable, ok := _self.rawParam["enable"]
	if !ok {
		return ok
	}
	var val int
	valueOf := reflect.ValueOf(enable)
	if valueOf.CanInt() {
		val = int(valueOf.Int())
	} else {
		kind := valueOf.Kind()
		switch kind {
		case reflect.Bool:
			return valueOf.Bool()
		case reflect.String:
			val, _ = strconv.Atoi(valueOf.String())
		default:
			val = 0
		}
	}
	return val == 1
}

func (_self commonLoginPolicyDTO) GetKey(key string) interface{} {
	val, ok := _self.rawParam[key]
	if ok {
		return val
	}
	return nil
}

func (_self commonLoginPolicyDTO) Verify() error {
	return nil
}

// OTPLoginPolicyDTO commonLoginPolicyDTO 的子类
// 如果一个子类对象可以用其父类对象替换，那么这个子类必须满足父类所有的接口。
// Go 语言的继承并不完全满足里氏替换原则
type OTPLoginPolicyDTO struct {
	// Go 语言没有显式的继承语法。在 Go 语言中，要实现继承，
	// 需要使用嵌入式结构体（embedded struct）和指针。
	// 嵌入式结构体允许在一个结构体中嵌入另一个结构体的字段，从而实现继承的效果。
	commonLoginPolicyDTO `json:"-"`
	// Issuer 签发名称
	Issuer string `json:"issuer" binding:"required,max=64"`
	// Period 间隔时长
	Period uint `json:"period" binding:"required,min=30,max=60"`
	// SecretSize
	SecretSize uint `json:"secretSize" binding:"required,min=12,max=24"`
}

func (_self OTPLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	_self.commonLoginPolicyDTO.rawParam = param
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &_self)
	return _self, err
}

func (_self OTPLoginPolicyDTO) Verify() error {
	return binding.Validator.ValidateStruct(_self)
}

type EmailLoginPolicyDTO struct {
	commonLoginPolicyDTO `json:"-"`
	ExpireTime           uint `json:"expireTime" binding:"required,min=2,max=10"`
}

func (_self EmailLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	var err error
	_self.commonLoginPolicyDTO.rawParam = param
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &_self)
	return _self, err
}

type GlobalLoginPolicyDTO struct {
	commonLoginPolicyDTO `json:"-"`
	Mode                 uint `json:"mode"  binding:"required,max=1,min=0"`
}

func (_self GlobalLoginPolicyDTO) New(param map[string]interface{}) (ILoginPolicyDTO, error) {
	var err error
	_self.commonLoginPolicyDTO.rawParam = param
	marshal, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &_self)
	return _self, err
}
