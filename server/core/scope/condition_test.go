package scope

import "testing"

func TestFormatLike(t *testing.T) {

	t.Log(formatLike("type"))
	t.Log(formatLikeRight("type"))
	t.Log(formatLikeLeft("type"))
}

func Test_formatLike(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1", args: args{arg: "arg1"}, want: "%arg1%"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatLike(tt.args.arg); got != tt.want {
				t.Errorf("formatLike() = %v, want %v", got, tt.want)
			}
		})
	}
}
