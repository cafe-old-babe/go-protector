package condition

import (
	"gorm.io/gorm/schema"
	"sync"
)

var parseMap sync.Map

var defaultStrategy schema.NamingStrategy

// GenerateCaseWhenSet 生成 set 语句
func GenerateCaseWhenSet(setColumn, whenColumn, thenColumn string,
	data interface{}) (setSql string, err error) {
	//if data == nil {
	//	err = c_error.ErrParamInvalid
	//	return
	//}
	//var dataSchema *schema.Schema
	//dataSchema, err = schema.Parse(data, &parseMap, defaultStrategy)
	//
	//var sqlBuilder strings.Builder
	//
	//for i, elem := range data {
	//	schema.ParseTagSetting()
	//	value := reflect.Indirect(reflect.ValueOf(elem))
	//	value.FieldByName(setColumn).Interface()
	//
	//}
	//setSql = sqlBuilder.String()
	return
}
