package utils

import (
	"fmt"
	"testing"
)

func TestGetRandomCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%s\n", GetRandomCode())
	}
}
