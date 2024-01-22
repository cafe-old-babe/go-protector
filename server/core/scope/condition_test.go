package scope

import "testing"

func TestFormatLike(t *testing.T) {

	t.Log(formatLike("type"))
	t.Log(formatLikeRight("type"))
	t.Log(formatLikeLeft("type"))
}
