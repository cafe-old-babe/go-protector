package cmd

type Handler interface {
	GetIndex() int
	GetId() uint64
	PassToClient(r rune) bool
	PassToServer(r rune) bool
}
