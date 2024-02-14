package c_translator

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

// ConvertValidateErr 转换校验错误
func ConvertValidateErr(err error) error {
	// github.com/gin-gonic/gin@v1.9.1/binding/binding.go:70
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {

		var messageSlice []string
		for _, v := range errs.Translate(trans) {
			messageSlice = append(messageSlice, v)
		}

		err = errors.New(strings.Join(messageSlice, ","))
	}
	return err

}

var trans ut.Translator

// https://github.com/go-playground/validator
// https://github.com/go-playground/validator/blob/master/_examples/translations/main.go

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			return name
		})
		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		trans, _ = uni.GetTranslator(zhT.Locale())
		_ = zhTranslations.RegisterDefaultTranslations(v, trans)
	}
	return
}
