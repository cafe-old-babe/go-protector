package c_jwt

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"go-protector/server/core/config"
	"go-protector/server/core/custom/c_error"
	"go-protector/server/core/utils"
	"go-protector/server/models/dto"
	time "time"
)

const secretKey = "secret-p90-23n09.32342"

// GenerateToken 生成token https://pkg.go.dev/github.com/golang-jwt/jwt/v5#Token
// go get -u github.com/golang-jwt/jwt/v5
func GenerateToken(currentUser *dto.CurrentUser) (jwtString string, expireAt time.Time, err error) {
	var bytes []byte
	sessionTimeout := config.GetConfig().Jwt.SessionTimeout
	if bytes, err = json.Marshal(currentUser); err != nil {
		return
	}
	if currentUser.SessionId <= 0 {
		currentUser.SessionId = utils.GetNextId()
	}
	now := time.Now().Local()
	expireAt = now.Add(time.Second * time.Duration(sessionTimeout))
	registeredClaims := jwt.RegisteredClaims{
		// 过期时间
		ExpiresAt: jwt.NewNumericDate(expireAt),
		// 签发时间
		IssuedAt: jwt.NewNumericDate(now),
		// 携带的内容
		Subject: string(bytes),
	}

	tokenPointer := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	jwtString, err = tokenPointer.SignedString([]byte(secretKey))
	return
}

// ParserToken 解析token
func ParserToken(jwtString *string) (userPointer *dto.CurrentUser, err error) {

	token, err := jwt.NewParser().ParseWithClaims(*jwtString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims := token.Claims
	// 获取currentUser
	subject, _ := claims.GetSubject()

	var currentUser dto.CurrentUser
	if err = json.Unmarshal([]byte(subject), &currentUser); err != nil {
		return
	}
	userPointer = &currentUser
	// 换token
	iat, _ := claims.GetIssuedAt()
	// 如果token的未过有效期, 但token时效性过期了,更换token
	tokenTimeout := config.GetConfig().Jwt.TokenTimeout
	if time.Now().After(iat.Add(time.Second * time.Duration(tokenTimeout))) {
		jwtString, err = ReGenerateToken(jwtString, userPointer)
	}
	return
}

// ReGenerateToken 重新生成Token,将之前的老token缓存一下
func ReGenerateToken(jwtString *string, currentUser *dto.CurrentUser) (newJwtTString *string, err error) {
	if jwtString == nil || len(*jwtString) <= 0 {
		err = c_error.ErrParamInvalid
		return
	}
	//redis := cache.GetRedis()
	// todo 并发换token
	*newJwtTString, _, err = GenerateToken(currentUser)
	return
}
