package c_jwt

import (
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"path/filepath"
	"testing"
	"time"
	_ "time/tzdata"
)

func init() {
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
}
func TestGenerateToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		abs, _ := filepath.Abs("/opt/work_space/github/go-protector/config/config.yml")
		t.Setenv(consts.EnvConfig, abs)
		jwtString, _, err := GenerateToken(&current.User{
			ID:        0,
			SessionId: "",
			LoginName: "",
			UserName:  "",
			LoginTime: "",
			LoginIp:   "",
		})
		if err != nil {
			t.Fatalf("generate fail: %v", err)
		}
		t.Logf("jwt: %v", jwtString)
	})
}

func TestParserToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		abs, _ := filepath.Abs("/opt/work_space/github/go-protector/config/config.yml")
		t.Setenv(consts.EnvConfig, abs)
		JwtStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJJRFwiOjAsXCJzZXNzaW9uSWRcIjowLFwibG9naW5OYW1lXCI6XCJcIixcInVzZXJOYW1lXCI6XCJcIixcImxvZ2luVGltZVwiOlwiXCIsXCJsb2dpbklwXCI6XCJcIn0iLCJleHAiOjE3MDQ4MTU1ODIsImlhdCI6MTcwNDgxNDk4Mn0.8zZdWDrzqdEWq6fAS9zBusvDI03m4HjvDEUrvF4vvIc"
		targetJwtStr := JwtStr

		currentUser, err := ParserToken(&targetJwtStr)
		if err != nil {
			t.Fatalf("parser err: %v", err)
		}
		if targetJwtStr != JwtStr {
			t.Logf("jwt modfiy")
		}
		t.Logf("currentUser: %v", *currentUser)
	})
}
