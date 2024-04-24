package c_jwt

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
	key := fmt.Sprintf(consts.OnlineUserCacheKeyFmt, currentUser.LoginName)
	err = redisClient.HSet(context.TODO(), key, currentUser.SessionId, jwtString).Err()
	return
}

// ParserToken 解析token
func ParserToken(jwtString *string) (userPointer *current.User, err error) {
	// todo  check redis --> token

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
	// 换token
	iat, _ := claims.GetIssuedAt()
	// 如果token的未过有效期, 但token时效性过期了,更换token
	tokenTimeout := config.GetConfig().Jwt.TokenTimeout
	if !iat.Add(time.Minute * time.Duration(tokenTimeout)).After(time.Now()) {
		jwtString, err = ReGenerateToken(jwtString, userPointer)

	}
	return
}

// ReGenerateToken 重新生成Token,将之前的老token缓存一下
func ReGenerateToken(jwtString *string, currentUser *current.User) (newJwtTString *string, err error) {
	if jwtString == nil || len(*jwtString) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	//redis := cache.GetRedis()
	// todo 并发换token
	newJwtTString, _, err = GenerateToken(currentUser)
	return
}
