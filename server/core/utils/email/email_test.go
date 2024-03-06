package email

import (
	"go-protector/server/core/config"
	"testing"
)

func TestSendImage(t *testing.T) {
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
					Email: &config.Email{
						Host:     "",
						Port:     0,
						Username: "",
						Password: "",
					},
					To:      "",
					Subject: "",
					Body:    "",
				},
				imageBase64: "",
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
				Email: &config.Email{
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
