package hdl

type Handler interface {
	Start(port int)
	Close() error
}
