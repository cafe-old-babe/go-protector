package api

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()
	panic("safasdfasdfasd")

}

func TestMap(t *testing.T) {
	testMap := map[string]string{
		"key": "value",
	}
	fmt.Printf("testMap: %v\n", testMap)
	testMap2 := map[string]interface{}{
		"key": 111,
	}
	fmt.Printf("testMap2: %v\n", testMap2)

}
