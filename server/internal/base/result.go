package base

import (
	"fmt"
	"go-protector/server/internal/custom/c_translator"
	"net/http"
)

type Result struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ResultSuccess 创建成功返回结果
// Result.Code 默认 http.StatusOK
// data 对应 Result.Data
// msg 对应 Result.Message, 默认操作成功,可以自定义,也可以格式化
// 2-12	【实战】封装统一的返回格式-掌握GO语言可变参数与切片的使用技巧
func ResultSuccess(data any, msg ...string) *Result {

	return ResultCustom(http.StatusOK, data, msg...)
}

func ResultSuccessMsg(msg ...string) *Result {

	return ResultCustom(http.StatusOK, nil, msg...)
}

func ResultFailure(data any, msg ...string) *Result {
	return ResultCustom(http.StatusBadRequest, data, msg...)
}
func ResultFailureMsg(msg ...string) *Result {
	return ResultCustom(http.StatusBadRequest, nil, msg...)
}

func ResultFailureErr(err error) *Result {
	err = c_translator.ConvertValidateErr(err)
	return ResultCustom(http.StatusBadRequest, nil, err.Error())
}

func ResultCustom(code int, data any, msg ...string) *Result {

	res := Result{}
	res.Code = code
	res.Data = data
	if len(msg) <= 0 {
		if code == http.StatusOK {
			res.Message = "操作成功"
		} else {
			res.Message = "操作失败"
		}
	} else if len(msg) == 1 {
		res.Message = msg[0]
	} else {
		res.Message = fmt.Sprintf(msg[0], msg[1:])
	}
	return &res
}

func ResultPage(data any, pagination IPagination, count int64) *Result {
	page := Page{
		TotalCount: int(count),
		PageNo:     pagination.GetPageIndex(),
		PageSize:   pagination.GetPageSize(),
		Data:       data,
	}
	if count > 0 {
		page.TotalPages = page.TotalCount / page.PageSize
		remainder := page.TotalCount % page.PageSize
		if remainder != 0 {
			page.TotalPages += 1
		}
	}
	return ResultSuccess(page, "查询成功")
}

// IsSuccess 判断是否成功
func (_self *Result) IsSuccess() bool {
	return _self.Code == http.StatusOK

}
