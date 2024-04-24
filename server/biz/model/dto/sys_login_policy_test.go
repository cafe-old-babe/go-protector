package dto

import (
	"fmt"
	"go-protector/server/internal/custom/c_translator"
	"reflect"
	"testing"
)

func TestGlobalLoginPolicyDTO_New(t *testing.T) {
	type fields struct {
		basePolicyDTO basePolicyDTO
		Mode          string
	}
	type args struct {
		param map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ILoginPolicyDTO
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				basePolicyDTO: basePolicyDTO{
					rawParam: nil,
					Enable:   "2",
				},
				Mode: "2",
			},
			args: args{
				param: map[string]interface{}{
					"mode":   "2",
					"enable": "2",
				},
			},
			want: &GlobalLoginPolicyDTO{
				basePolicyDTO: basePolicyDTO{
					rawParam: map[string]interface{}{
						"mode":   "2",
						"enable": "2",
					},
					Enable: "2",
				},
				Mode: "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_self := GlobalLoginPolicyDTO{
				basePolicyDTO: tt.fields.basePolicyDTO,
				Mode:          tt.fields.Mode,
			}
			got, err := _self.New(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
			if err := _self.Validate(got); (err != nil) != tt.wantErr {
				err = c_translator.ConvertValidateErr(err)
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMap(t *testing.T) {
	target := make(map[string]string)
	targetSlice := make([]string, 0)
	modifyMap(target)
	modifySlice(&targetSlice)
	fmt.Printf("%v\n", target)
	fmt.Printf("%v\n", targetSlice)
}

func modifySlice(slice *[]string) {
	*slice = append(*slice, "1")
}

func modifyMap(target map[string]string) {
	target["1"] = "1"
}
