package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"go-protector/server/core/consts"
	"go-protector/server/core/initialize"
	"go-protector/server/models/dto"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var configPath string

func init() {
	configPath, _ = filepath.Abs("/opt/work_space/github/go-protector/config/config.yml")
	_ = os.Setenv(consts.EnvConfig, configPath)
}

func TestServer(t *testing.T) {
	var server *http.Server
	var err error
	if server, err = initialize.GetServer(); err != nil {
		t.Fatalf("get server err: %v", err)
	}
	t.Run("setStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		var buffer bytes.Buffer
		login := dto.SetStatus{
			ID:         1,
			UserStatus: 0,
		}
		marshal, _ := json.Marshal(login)
		buffer.Write(marshal)
		req, _ := http.NewRequest("POST", "/api/user/setStatus", &buffer)

		server.Handler.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		t.Logf("response: %v\n", w.Body.String())
	})

	t.Run("login/success", func(t *testing.T) {
		w := httptest.NewRecorder()
		var buffer bytes.Buffer
		login := dto.Login{
			LoginName: "admin",
			Password:  "888888",
		}
		marshal, _ := json.Marshal(login)
		buffer.Write(marshal)
		req, _ := http.NewRequest("POST", "/api/user/login", &buffer)
		req.Header.Set("X-Real-IP", "172.16.1.1")
		server.Handler.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		t.Logf("response: %v\n", w.Body.String())
	})

	t.Run("login/failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		var buffer bytes.Buffer
		login := dto.Login{
			LoginName: "admin",
			Password:  "8888880",
		}
		marshal, _ := json.Marshal(login)
		buffer.Write(marshal)
		req, _ := http.NewRequest("POST", "/api/user/login", &buffer)

		server.Handler.ServeHTTP(w, req)
		assert.Equal(t, w.Code, http.StatusOK)
		t.Logf("response: %v\n", w.Body.String())
	})

}
