package pipeline

import "github.com/nats-io/nats.go"

func NewDriverFactory[T any](natsConn *nats.Conn) DriverFactory[T] {
	return DriverFactory[T]{
		natsConn: natsConn,
	}
}

type DriverFactory[T any] struct {
	natsConn *nats.Conn
}

func (df DriverFactory[T]) NewDriver(channelName string) WorkerDriver[T] {
	if df.natsConn != nil {
		return NewNatsDriver[T](df.natsConn, channelName)
	}

	return NewInMemoryQueueDriver[T](channelName)
}
