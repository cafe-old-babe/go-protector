package c_captcha

import (
	"context"
	"github.com/mojocn/base64Captcha"
	"go-protector/server/core/cache"
	"go-protector/server/core/consts"
	"time"
)

var store redisStore

// Generate 生成图片
func Generate() (id, b64s string, err error) {
	c := base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, store)

	id, b64s, _, err = c.Generate()
	return
}

// Verify 校验
func Verify(id, answer string, clear bool) bool {
	return store.Verify(id, answer, clear)
}

type redisStore struct {
}

const CaptchaPrefix = consts.CachePrefix + ":captcha:"

// Set 设置
func (r redisStore) Set(id string, value string) error {
	//有效时间10分钟
	cache.GetRedisClient().Set(context.TODO(), CaptchaPrefix+id, value, time.Minute*10)
	return nil
}

// Get 获取
func (r redisStore) Get(id string, clear bool) string {
	key := CaptchaPrefix + id

	redisCli := cache.GetRedisClient()
	val, err := redisCli.Get(context.TODO(), key).Result()

	if err != nil {
		return ""
	}
	if clear {
		//clear为true，验证通过，删除这个验证码
		redisCli.Del(context.TODO(), key)
	}
	return val
}

// Verify 验证
func (r redisStore) Verify(id, answer string, clear bool) bool {
	if id == "" || answer == "" {
		return false
	}
	return r.Get(id, clear) == answer
}
