package cmd

type Handler interface {
	GetIndex() int
	GetId() uint64
	PassToClient(r rune)
	PassToServer(r rune) bool
	Close()
}
