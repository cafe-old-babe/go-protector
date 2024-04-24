package gm

import (
	"encoding/base64"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
	"go-protector/server/internal/config"
	"sync"
)

var keyByte []byte
var once sync.Once

// Sm4EncryptCBC sms4 CBC 加密
func Sm4EncryptCBC(deStr string) (enStr string, err error) {
	once.Do(func() {
		keyByte = []byte(config.GetConfig().Server.Sm4Key)
	})
	var out []byte
	if out, err = sm4.Sm4Cbc(keyByte, []byte(deStr), true); err != nil {
		return
	}
	enStr = base64.StdEncoding.EncodeToString(out)
	return
}

// Sm4DecryptCBC sms4 CBC 解密
func Sm4DecryptCBC(enStr string) (deStr string, err error) {

	once.Do(func() {
		keyByte = []byte(config.GetConfig().Server.Sm4Key)
	})
	var bytes []byte
	if bytes, err = base64.StdEncoding.DecodeString(enStr); err != nil {
		return
	}
	var out []byte
	if out, err = sm4.Sm4Cbc(keyByte, bytes, false); err != nil {
		return
	}

	deStr = fmt.Sprintf("%s", out)
	return
}
