package dto

type Page struct {
	Count      int `json:"totalCount"`
	PageIndex  int `json:"pageNo"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPage"`
	List       any `json:"list"`
}

type IPagination interface {
	GetPageIndex() int
	GetPageSize() int
	GetPagination() IPagination
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

func (_self *Pagination) GetPagination() IPagination {
	return _self
}
