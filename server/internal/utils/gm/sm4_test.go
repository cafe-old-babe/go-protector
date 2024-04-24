package gm

import (
	"encoding/json"
	"fmt"
	"go-protector/server/internal/consts"
	"os"
	"path/filepath"
	"testing"
)

var configPath string

func init() {
	configPath, _ = filepath.Abs("/opt/work_space/github/go-protector/config/config.yml")
	_ = os.Setenv(consts.EnvConfig, configPath)
}

type args struct {
	DeStr string `json:"deStr"`
}

func TestSm4EncryptCBC(t *testing.T) {

	arg := args{DeStr: "123456"}
	marshal, _ := json.Marshal(arg)

	str, _ := Sm4EncryptCBC(string(marshal))
	fmt.Printf("en: %v\n", str)
	/*
		tests := []struct {
			name      string
			args      args
			wantEnStr string
			wantErr   bool
		}{
			{
				name:      "1",
				args:      args{DeStr: string(marshal)},
				wantEnStr: "s2y6dbBFus/h6q8crE/TPw==",
				wantErr:   true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotEnStr, err := Sm4EncryptCBC(tt.args.DeStr)
				if (err != nil) != tt.wantErr {
					t.Errorf("Sm4EncryptCBC() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotEnStr != tt.wantEnStr {
					t.Errorf("Sm4EncryptCBC() gotEnStr = %v, want %v", gotEnStr, tt.wantEnStr)
				}
				t.Logf("en: %s\n", gotEnStr)

			})
		}*/
}

func TestSm4DecryptCBC(t *testing.T) {

	str, _ := Sm4DecryptCBC("S65TqUfgcN1ZXPQ/+be6NB/Cf2OrzPavtaedNAv/7XM=")
	fmt.Printf("de: %v\n", str)
	var arg args

	_ = json.Unmarshal([]byte(str), &arg)

	fmt.Println(arg.DeStr)

	/*type args struct {
		enStr string
	}
	tests := []struct {
		name      string
		args      args
		wantDeStr string
		wantErr   bool
	}{
		{
			name:      "1",
			args:      args{enStr: "s2y6dbBFus/h6q8crE/TPw=="},
			wantDeStr: "123456",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeStr, err := Sm4DecryptCBC(tt.args.enStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sm4DecryptCBC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDeStr != tt.wantDeStr {
				t.Errorf("Sm4DecryptCBC() gotDeStr = %v, want %v", gotDeStr, tt.wantDeStr)
			}
		})
	}*/
}
