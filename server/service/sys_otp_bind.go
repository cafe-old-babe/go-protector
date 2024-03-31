package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go-protector/server/core/base"
	"go-protector/server/core/current"
	"go-protector/server/models/dto"
	"go-protector/server/models/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"image/png"
	"time"
)

type SysOtpBind struct {
	base.Service
}

// GetQRCodeBase64ByUser go get github.com/pquerna/otp
func (_self *SysOtpBind) GetQRCodeBase64ByUser(user current.User,
	policyDTO dto.OTPLoginPolicyDTO) (imageBase64 string, err error) {
	var model entity.SysOtpBind
	if err = _self.DB.First(&model, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if model, err = _self.CreateOtpByUser(user, policyDTO); err != nil {
				return
			}
		}
		return
	}
	var key *otp.Key
	key, err = otp.NewKeyFromURL(model.OtpAuthURL)
	if err != nil {
		return
	}
	if key.Issuer() != policyDTO.Issuer ||
		key.Period() != uint64(policyDTO.Period) ||
		len(key.Secret()) != int(policyDTO.SecretSize) {
		if model, err = _self.CreateOtpByUser(user, policyDTO); err != nil {
			return
		}
		key, err = otp.NewKeyFromURL(model.OtpAuthURL)
		if err != nil {
			return
		}
	}

	image, err := key.Image(200, 200)
	imageBuffer := new(bytes.Buffer)

	err = png.Encode(imageBuffer, image)
	if err != nil {
		_self.Logger.Error("无法保存二维码: %v", err)
		return
	}
	imageBase64 = base64.StdEncoding.EncodeToString(imageBuffer.Bytes())
	return
}

func (_self *SysOtpBind) VerifyCode(user current.User, code string) (err error) {
	var model entity.SysOtpBind
	if err = _self.DB.First(&model, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("动态认证码无效或已过期")
		}
		return
	}
	var key *otp.Key
	key, err = otp.NewKeyFromURL(model.OtpAuthURL)
	now := time.Now()
	var pass bool
	if pass, err = totp.ValidateCustom(code, key.Secret(), now, totp.ValidateOpts{
		Period:    uint(key.Period()),
		Skew:      1,
		Digits:    key.Digits(),
		Algorithm: key.Algorithm(),
	}); err != nil {
		return
	}
	if !pass {
		err = errors.New("动态认证码无效或已过期")
	}
	return
}

func (_self *SysOtpBind) CreateOtpByUser(user current.User,
	policyDTO dto.OTPLoginPolicyDTO) (model entity.SysOtpBind, err error) {

	if err = policyDTO.Validate(&policyDTO); err != nil {
		return entity.SysOtpBind{}, err
	}

	// 生成一个新的密钥
	var key *otp.Key
	if key, err = totp.Generate(totp.GenerateOpts{
		Issuer:      policyDTO.Issuer,
		AccountName: user.LoginName,
		Period:      policyDTO.Period,
		SecretSize:  policyDTO.SecretSize,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	}); err != nil {
		return
	}
	model = entity.SysOtpBind{
		UserId:     user.ID,
		OtpAuthURL: key.URL(),
	}
	err = _self.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&model).Error
	return
}
