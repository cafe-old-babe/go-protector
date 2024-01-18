package dto

import (
	"fmt"
	"net/http"
)

var src Result

type Result struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ResultSuccess 创建成功返回结果
// Result.Code 默认 http.StatusOK
// data 对应 Result.Data
// msg 对应 Result.Message, 默认操作成功,可以自定义,也可以格式化
func ResultSuccess(data any, msg ...string) *Result {

	return ResultCustom(http.StatusOK, data, msg...)
}

func ResultFailure(data any, msg ...string) *Result {
	return ResultCustom(http.StatusBadRequest, data, msg...)
}
func ResultFailureMsg(msg ...string) *Result {
	return ResultCustom(http.StatusBadRequest, nil, msg...)
}

func ResultFailureErr(err error) *Result {
	return ResultCustom(http.StatusBadRequest, nil, err.Error())
}

func ResultCustom(code int, data any, msg ...string) *Result {
	// res := &src 和 new(Result) 的效率是相同的。
	// 使用 res := src 的方式是可行的，但是它不如使用 res := &src 的方式高效。
	//
	// 使用 := 赋值操作符将一个变量赋值给另一个变量时，实际上是在内存中创建一个新的变量，
	// 并将原变量的值复制到新变量中。因此，使用 res := src 的方式需要先创建一个指向 Result 结构体的指针 src，
	// 然后再将其赋值给 res 变量，这样就多了一次内存分配和指针拷贝的过程，效率较低。
	//
	// 而使用 res := &src 的方式可以直接将 src 变量的地址赋值给 res 变量，
	// 避免了指针拷贝的过程，效率更高。同时，使用 res := &src 的方式也更清晰、易读，
	// 可以清楚地表示 res 变量是一个指向 src 变量的指针。
	res := &src
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
	return res
}

type Page struct {
	Count      int64 `json:"totalCount"`
	PageIndex  int   `json:"pageNo"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPage"`
	List       any   `json:"list"`
}

type Pagination struct {
	PageIndex int `json:"current" form:"current"`
	PageSize  int `json:"pageSize" form:"pageSize"`
}

func (_self *Pagination) GetPageIndex() int {
	if _self.PageIndex <= 0 {
		_self.PageIndex = 1
	}
	return _self.PageIndex
}

func (_self *Pagination) GetPageSize() int {
	if _self.PageSize <= 0 {
		_self.PageSize = 10
	}
	return _self.PageSize
}
