package api

import "testing"

func TestError(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	panic("safasdfasdfasd")

}
