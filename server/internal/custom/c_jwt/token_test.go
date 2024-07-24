package c_jwt

import (
	"bytes"
	"fmt"
	"go-protector/server/internal/consts"
	"go-protector/server/internal/current"
	"net/http"
	"path/filepath"
	"sync"
	"sync/atomic"
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

// 3-14	【扩展】本章小结，编写并发续约测试用例-掌握同步原语sync.WaitGroup
func TestConcurrentToken(t *testing.T) {

	var token = "JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ7XCJpZFwiOjEsXCJzZXNzaW9uSWRcIjpcIjU3MjI0MzEwMzIyMjU5NTU4NFwiLFwibG9naW5OYW1lXCI6XCJhZG1pblwiLFwiZW1haWxcIjpcImNhZmUtb2xkLWJhYmVAcXEuY29tXCIsXCJ1c2VyTmFtZVwiOlwi566h55CG5ZGYXCIsXCJsb2dpblRpbWVcIjpcIjIwMjQtMDQtMjggMDI6MDk6NDdcIixcImxvZ2luSXBcIjpcIjEyNy4wLjAuMVwiLFwiYXZhdGFyXCI6XCJodHRwczovL2d3LmFsaXBheW9iamVjdHMuY29tL3pvcy9ybXNwb3J0YWwvQmlhemZhbnhtYW1OUm94eFZ4a2EucG5nXCIsXCJyb2xlSWRzXCI6WzFdLFwiZGVwdElkXCI6MixcImlzQWRtaW5cIjp0cnVlfSIsImV4cCI6MTcxNDI0NTg2MiwiaWF0IjoxNzE0MjQ0MDYyfQ.9z6tyTJ2ugTHN8y6G5nWI686t65wB8Q1GNixvLR7gbk"
	// 指定 URL
	url := "http://127.0.0.1:8888/api/asset-gateway/page"

	// 创建一个 WaitGroup 实例
	var wg sync.WaitGroup

	// 创建一个 channel 用于接收响应结果
	var count atomic.Int32

	// 启动多个 goroutine 进行并发测试
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()

			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{"pageNo":1,"pageSize":10}`)))
			if err != nil {
				fmt.Errorf("创建请求失败: %v\n", err)
				return
			}
			req.Header.Set("Authorization", token)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Errorf("请求失败: %v\n", err)
				return
			}
			defer resp.Body.Close()

			//body, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Println("读取响应结果失败:", err)
				return
			}
			authorization := resp.Header.Get("Authorization")
			if authorization != "" {
				fmt.Printf("%d--> count: %d, 获取到 jwt header: %s\n", j, count.Add(1), authorization)

			}
		}(i)
	}

	wg.Wait()

}
