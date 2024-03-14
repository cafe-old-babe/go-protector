package dto

import (
	"encoding/json"
	"fmt"
	"go-protector/server/core/consts"
	"go-protector/server/core/custom/c_translator"
	"go-protector/server/core/custom/c_type"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	type args struct {
		t      c_type.LoginPolicyCode
		enable int
		PJson  string
	}
	otpWant := OTPLoginPolicyDTO{
		Issuer:     "hhh",
		Period:     30,
		SecretSize: 0,
	}
	marshal, _ := json.Marshal(otpWant)
	_ = json.Unmarshal(marshal, &otpWant.rawParam)

	tests := []struct {
		name    string
		args    args
		want    ILoginPolicyDTO
		wantErr bool
	}{
		{
			name: "otp",
			args: args{
				t:     consts.LoginPolicyOtp,
				PJson: string(marshal),
			},
			want:    otpWant,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateLoginPolicyDTO(tt.args.t, tt.args.PJson)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLoginPolicyDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateLoginPolicyDTO() got = %v, want %v", got, tt.want)
			}
			if err = got.Verify(); err != nil {
				err = c_translator.ConvertValidateErr(err)
				t.Errorf("Verify error = %v", err)
			}
			fmt.Printf("%v\n", got)
		})
	}
}
