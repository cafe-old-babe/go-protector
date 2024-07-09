package observe

type Subject interface {
	RegisterObserver(observer Observer) error
	RemoveObserver(observer Observer)
	NotifyUpdateObservers(id uint64, r rune)
	NotifyCloseObservers(id uint64)
}
