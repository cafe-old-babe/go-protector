package observe

type Observer interface {
	GetSsoId() uint64
	GetObId() uint64
	Update(rune)
	Close()
}
