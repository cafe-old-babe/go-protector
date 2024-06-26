package cmd

var DefaultCmdHandler defaultCmdHandler

type defaultCmdHandler struct {
}

func (_self *defaultCmdHandler) GetIndex() int {
	return -9999
}

func (_self *defaultCmdHandler) PassToClient(r rune) bool {
	return true
}

func (_self *defaultCmdHandler) PassToServer(r rune) bool {
	return true
}
