package backup

type Creator interface {
	Create() error
}
