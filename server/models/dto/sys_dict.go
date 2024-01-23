package dto

type DictTypePageReq struct {
	Pagination
	TypeCode string `json:"typeCode"`
	TypeName string `json:"typeName"`
}

type DictDataPageReq struct {
	DictTypePageReq
	DataCode string `json:"dataCode"`
	DataName string `json:"dataName"`
}

type DictDataUpdateStatusReq struct {
	ID     string `uri:"id"`
	Status int    `uri:"status"`
}
