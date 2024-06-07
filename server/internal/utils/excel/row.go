package excel

type ColInfo struct {
	Name  string
	Title string
	Index int
	Width int
}

type Row interface {
	GetLineNum() int
	SetLineNum(int)
	SetErr(error)
	GetErr() string
}

type StdRow struct {
	LineNum int
	ErrMsg  string `excel:"title:错误消息;width:50;index:999"`
}

func (_self *StdRow) GetLineNum() int {
	return _self.LineNum
}

func (_self *StdRow) SetLineNum(i int) {
	_self.LineNum = i
}

func (_self *StdRow) SetErr(err error) {
	_self.ErrMsg = err.Error()
}

func (_self *StdRow) GetErr() string {
	return _self.ErrMsg
}
