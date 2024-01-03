package utils

import "github.com/GUAIK-ORG/go-snowflake/snowflake"

var _innerSnowflake *snowflake.Snowflake

func init() {
	var err error
	_innerSnowflake, err = snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		panic(err.Error())
	}
}
func GetNextId() (nextId uint64) {
	nextId = uint64(_innerSnowflake.NextVal())
	return
}
