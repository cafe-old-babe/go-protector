package vo

type DictDataPage struct {
	ID             uint64 `json:"id"`
	TypeCode       string `json:"typeCode"`
	TypeName       string `json:"typeName"`
	DataCode       string `json:"dataCode"`
	DataName       string `json:"dataName"`
	Sort           int    `json:"sort"`
	DataStatus     string `json:"dataStatus"`
	DataStatusText string `json:"dataStatusText"`
}
