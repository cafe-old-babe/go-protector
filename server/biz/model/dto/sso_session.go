package dto

type ConnectBySessionReq struct {
	Id uint64
	H  int `form:"h"`
	W  int `form:"w"`
}
