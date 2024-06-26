package cmd

import (
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	var i int
	i = strings.IndexAny("\"`'", "`'\"")
	println(i)
	i = strings.IndexAny("`\"'", "`'\"")
	println(i)
	i = strings.IndexAny("`'\"", "`'\"")
	println(i)
	i = strings.IndexAny("'\"`", "`'\"")
	println(i)
	i = strings.IndexAny("\"`'", "`'\"")
	println(i)
}
