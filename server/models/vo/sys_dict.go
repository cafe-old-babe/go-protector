package vo

type DictDataPage struct {
	ID         uint64 `json:"id"`
	TypeCode   string `json:"typeCode"`
	TypeName   string `json:"typeName"`
	DataCode   string `json:"dataCode"`
	DataName   string `json:"dataName"`
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}

type DictDataList struct {
	ID       uint64 `json:"id"`
	DataCode string `json:"code"`
	DataName string `json:"text"`
}
