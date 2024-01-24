package dto

type IdsReq struct {
	Ids []uint64 `json:"ids"`
	ID  uint64   `json:"id"`
}

// GetIds 获取ids
func (_self IdsReq) GetIds() []uint64 {
	if _self.ID > 0 {
		_self.Ids = append(_self.Ids, _self.ID)
	}
	return _self.Ids
}
