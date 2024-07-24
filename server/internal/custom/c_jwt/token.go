package c_jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go-protector/server/internal/cache"
	"go-protector/server/internal/config"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"go-protector/server/internal/custom/c_error"
	"go-protector/server/internal/utils"
	"golang.org/x/net/context"
	"time"
)

const secretKey = "secret-p90-23n09.32342"

// GenerateToken 生成token https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Token
// go get -u github.com/golang-jwt/jwt/v5
// 3-9	【实战】登录成功后生成认证Token-掌握JWT原理及应用-掌握解决基于Session认证的局限性，RFC文档介绍，设置本地时区
func GenerateToken(currentUser *current.User) (jwtStringPointer *string, expireAt time.Time, err error) {
	var bytes []byte
	sessionTimeout := config.GetConfig().Jwt.SessionTimeout
	if len(currentUser.SessionId) <= 0 {
		currentUser.SessionId = utils.GetNextIdStr()
	}
	if bytes, err = json.Marshal(currentUser); err != nil {
		return
	}
	now := time.Now().Local()
	duration := time.Minute * time.Duration(sessionTimeout)
	expireAt = now.Add(duration)
	registeredClaims := jwt.RegisteredClaims{
		// 过期时间
		ExpiresAt: jwt.NewNumericDate(expireAt),
		// 签发时间
		IssuedAt: jwt.NewNumericDate(now),
		// 携带的内容
		Subject: string(bytes),
	}

	tokenPointer := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	var jwtString string
	if jwtString, err = tokenPointer.SignedString([]byte(secretKey)); err != nil {
		return
	}
	jwtStringPointer = &jwtString
	redisClient := cache.GetRedisClient()
	key := fmt.Sprintf(consts.OnlineUserCacheKeyFmt, currentUser.LoginName, currentUser.SessionId)
	err = redisClient.Set(context.TODO(), key, jwtString, duration).Err()
	return
}

// ParserToken 解析token
func ParserToken(jwtString *string) (userPointer *current.User, err error) {

	token, err := jwt.NewParser().ParseWithClaims(*jwtString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims := token.Claims
	// 获取currentUser
	subject, _ := claims.GetSubject()

	var currentUser current.User
	if err = json.Unmarshal([]byte(subject), &currentUser); err != nil {
		return
	}
	userPointer = &currentUser
	var breakCheckTimeout bool
	if breakCheckTimeout, err = DoCheckTokenEffective(jwtString, userPointer); err != nil {
		return
	}
	if breakCheckTimeout {
		return
	}
	// 换token
	iat, _ := claims.GetIssuedAt()
	// 如果token的未过有效期, 但token时效性过期了,更换token
	tokenTimeout := config.GetConfig().Jwt.TokenTimeout
	if !iat.Add(time.Minute * time.Duration(tokenTimeout)).After(time.Now()) {
		var newJwtString *string
		if newJwtString, err = ReGenerateToken(jwtString, userPointer); err == nil {
			*jwtString = *newJwtString
		}
	}
	return
}

// DoCheckTokenEffective check redis --> token
// 3-10	【实战】JWT主动销毁及续约-掌握并发场景下续约的两种解决方案
func DoCheckTokenEffective(jwtToken *string, currentUser *current.User) (breakCheckTimeout bool, err error) {
	keyFmt := consts.OnlineUserCacheKeyFmt

	redisClient := cache.GetRedisClient()
	var key, redisJwtToken, currentJwtToken string

	for {
		key = fmt.Sprintf(keyFmt, currentUser.LoginName, currentUser.SessionId)
		redisJwtToken, err = redisClient.Get(context.Background(), key).Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				return
			}
		}
		if redisJwtToken == *jwtToken {
			// effective
			err = nil
			// 方案二:最大限度的通知前端更换token, 如果不修改, 只会有一个线程(续约的线程)会通知到前台
			if keyFmt == consts.OnlineUserCacheLastKeyFmt {
				*jwtToken = currentJwtToken
				breakCheckTimeout = true
			}
			return
		}

		if keyFmt != consts.OnlineUserCacheLastKeyFmt {
			keyFmt = consts.OnlineUserCacheLastKeyFmt
			currentJwtToken = redisJwtToken
			_ = []byte(currentJwtToken)
			continue
		}
		err = c_error.ErrAuthFailure
		return

	}

}

// ReGenerateToken 重新生成Token,将之前的老token缓存一下
func ReGenerateToken(jwtString *string, currentUser *current.User) (newJwtTString *string, err error) {
	if jwtString == nil || len(*jwtString) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	redisClient := cache.GetRedisClient()
	key := fmt.Sprintf(consts.OnlineUserCacheLastKeyFmt, currentUser.LoginName, currentUser.SessionId)
	var setStatus bool
	// 并发换token
	if setStatus, err = redisClient.SetNX(context.Background(), key, jwtString, time.Minute).Result(); err != nil {
		return
	}
	if setStatus {
		// 方案一: 只有一个线程才可以续约token
		newJwtTString, _, err = GenerateToken(currentUser)
	} else {
		newJwtTString = jwtString
	}

	return
}
