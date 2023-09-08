package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jdvr/go-again"
)

func NewPostgresQueueDriver[T any](pool *pgxpool.Pool, channelName string) *postgresQueueDriver[T] {
	return &postgresQueueDriver[T]{
		channelName: channelName,
		pool:        pool,
		channels:    map[string]*channel[T]{},
		exit:        make(chan bool),
	}
}

// postgresQueueDriver is a queue driver that uses Postgres LISTEN/NOTIFY
// Since each queue needs its own connection, it's not practical/scalable
// to create a new Driver instance for each queue. Instead, we create a
// single Driver instance and use it to create channels for each queue.
//
// This driver requires one connection that listens to messages in any queue
// and routes them to the correct worker.
type postgresQueueDriver[T any] struct {
	channelName string
	pool        *pgxpool.Pool
	channels    map[string]*channel[T]
	running     bool
	exit        chan bool
}

func (qd *postgresQueueDriver[T]) getChannel(name string) (*channel[T], error) {
	ch, ok := qd.channels[name]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", name)
	}

	return ch, nil
}

type pgJob[T any] struct {
	Channel string `json:"channel"`
	Item    T      `json:"job"`
}

func (qd *postgresQueueDriver[T]) Start() {
	if qd.running {
		// we want only 1 worker here
		return
	}
	qd.running = true

	go func(qd *postgresQueueDriver[T]) {

		conn, err := qd.pool.Acquire(context.Background())
		if err != nil {
			panic(fmt.Errorf("error acquiring connection: %w", err))
		}
		defer conn.Release()

		for {
			select {
			case <-qd.exit:
				return
			default:
				qd.worker(conn)
			}
		}
	}(qd)
}

func (qd *postgresQueueDriver[T]) worker(conn *pgxpool.Conn) {
	_, err := conn.Exec(context.Background(), "listen "+qd.channelName)
	if err != nil {
		return
	}
	notification, err := conn.Conn().WaitForNotification(context.Background())
	if err != nil {
		return
	}

	job := pgJob[T]{}
	err = json.Unmarshal([]byte(notification.Payload), &job)
	if err != nil {
		return
	}

	channel, err := qd.getChannel(job.Channel)
	if err != nil {
		return
	}

	// spin off so we can keep listening for jobs
	channel.listener.Listen(job.Item)
}

func (qd *postgresQueueDriver[T]) Stop() {
	qd.exit <- true
}

// Channel registers a new queue channel and returns it
func (qd *postgresQueueDriver[T]) Channel(name string) *channel[T] {
	if _, channelNameExists := qd.channels[name]; channelNameExists {
		panic(fmt.Errorf("channel %s already exists", name))
	}

	ch := &channel[T]{
		postgresQueueDriver: qd,
		pool:                qd.pool,
		name:                name,
		enqueueTimeout:      defaultEnqueueTimeout,
	}

	qd.channels[name] = ch

	return ch
}

const defaultEnqueueTimeout = 5 * time.Minute

type channel[T any] struct {
	*postgresQueueDriver[T]
	name           string
	enqueueTimeout time.Duration
	pool           *pgxpool.Pool
	listener       Listener[T]
}

func (ch *channel[T]) SetListener(l Listener[T]) {
	ch.listener = l
}

func (ch *channel[T]) Enqueue(item T) {

	jj, err := json.Marshal(pgJob[T]{
		Channel: ch.name,
		Item:    item,
	})

	if err != nil {
		return
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), ch.enqueueTimeout)
	defer cancelCtx()

	conn, err := again.Retry[*pgxpool.Conn](ctx, func(ctx context.Context) (*pgxpool.Conn, error) {
		return ch.pool.Acquire(context.Background())
	})

	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Query(ctx, fmt.Sprintf(`select pg_notify('%s', $1)`, ch.name), jj)
	if err != nil {
		return
	}

}
