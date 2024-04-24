package email

import (
	"go-protector/server/internal/config"
	"go-protector/server/internal/consts"
	"path/filepath"
	"testing"
)

func TestSendImage(t *testing.T) {
	abs, _ := filepath.Abs("/opt/work_space/github/go-protector/config/config.yml")
	t.Setenv(consts.EnvConfig, abs)
	type args struct {
		dto         SendDTO
		imageBase64 string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				dto: SendDTO{
					To:      "cafe-old-babe@qq.com",
					Subject: "343434",
					Body:    "121321",
				},
				imageBase64: "data:image/png;base64,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",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendImage(tt.args.dto, tt.args.imageBase64); (err != nil) != tt.wantErr {
				t.Errorf("SendImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSend(t *testing.T) {
	type args struct {
		dto *SendDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{dto: &SendDTO{
				Email: config.Email{
					Host:     "",
					Port:     0,
					Username: "",
					Password: "",
				},
				To:      "",
				Subject: "",
				Body:    "",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Send(tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
