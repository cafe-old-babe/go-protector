package dto

type Page struct {
	TotalCount int `json:"totalCount"`
	PageNo     int `json:"pageNo"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPage"`
	Data       any `json:"data"`
}

type IPagination interface {
	GetPageIndex() int
	GetPageSize() int
	GetPagination() IPagination
}

type Pagination struct {
	PageNo   int `json:"pageNo" form:"pageNo"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func (_self *Pagination) GetPageIndex() int {
	if _self.PageNo <= 0 {
		_self.PageNo = 1
	}
	return _self.PageNo
}

func (_self *Pagination) GetPageSize() int {
	if _self.PageSize <= 0 {
		_self.PageSize = 10
	}
	return _self.PageSize
}

func (_self *Pagination) GetPagination() IPagination {
	return _self
}
