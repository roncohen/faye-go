package protocol

type Connection interface {
	Send([]Message) error
	IsConnected() bool
	IsSingleShot() bool
	Close()
	Priority() int
}
