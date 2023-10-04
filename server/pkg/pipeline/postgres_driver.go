package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jdvr/go-again"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

func NewPostgresQueueDriver[T any](pool *pgxpool.Pool, channelName string) *postgresQueueDriver[T] {
	id := id.GenerateID()
	return &postgresQueueDriver[T]{
		log:         newLoggerFn("PostgresQueueDriver - " + id.String()),
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
	log         loggerFn
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
		qd.log("already running")
		return
	}
	qd.running = true

	go func(qd *postgresQueueDriver[T]) {
		qd.log("start")

		for {
			select {
			case <-qd.exit:
				qd.log("exit")
				return
			default:
				_, err := again.Retry(context.Background(), func(_ context.Context) (bool, error) {
					qd.log("acquiring connection")
					conn, err := qd.pool.Acquire(context.Background())
					if err != nil {
						err = fmt.Errorf("error acquiring connection: %w", err)
						qd.log("%s", err.Error())
						return false, err
					}
					defer conn.Release()

					err = qd.worker(conn)
					if err != nil {
						err = fmt.Errorf("error in worker: %w", err)
						qd.log("%s", err.Error())
						return false, err
					}
					return true, nil
				})
				if err != nil {
					// this panic is intentional. forces the app to crash and restart
					panic(err)
				}
			}
		}
	}(qd)
}

func (qd *postgresQueueDriver[T]) worker(conn *pgxpool.Conn) error {
	qd.log("listening for notifications")
	_, err := conn.Exec(context.Background(), "listen "+qd.channelName)
	if err != nil {
		return fmt.Errorf("error listening for notifications: %w", err)
	}
	qd.log("waiting for notification")
	notification, err := conn.Conn().WaitForNotification(context.Background())
	if err != nil {
		return fmt.Errorf("error waiting for notification: %w", err)
	}

	job := pgJob[T]{}
	err = json.Unmarshal([]byte(notification.Payload), &job)
	if err != nil {
		// this error is not fatal. we can ignore it and continue
		qd.log("error unmarshalling pgJob: %s", err.Error())
		return nil
	}

	qd.log("received job for channel: %s", job.Channel)

	channel, err := qd.getChannel(job.Channel)
	if err != nil {
		// this error is not fatal. we can ignore it and continue
		qd.log("error getting channel: %s", err.Error())
		return nil
	}

	qd.log("processing job for channel: %s", job.Channel)
	channel.listener.Listen(job.Item)
	return nil
}

func (qd *postgresQueueDriver[T]) Stop() {
	qd.log("stopping")
	qd.exit <- true
}

// Channel registers a new queue channel and returns it
func (qd *postgresQueueDriver[T]) Channel(name string) *channel[T] {
	if _, channelNameExists := qd.channels[name]; channelNameExists {
		panic(fmt.Errorf("channel %s already exists", name))
	}

	ch := &channel[T]{
		postgresQueueDriver: qd,
		name:                name,
		log:                 newLoggerFn(fmt.Sprintf("PostgresQueueDriver - %s", name)),
		pool:                qd.pool,
	}

	qd.channels[name] = ch

	return ch
}

type channel[T any] struct {
	*postgresQueueDriver[T]
	name     string
	log      loggerFn
	pool     *pgxpool.Pool
	listener Listener[T]
}

func (ch *channel[T]) SetListener(l Listener[T]) {
	ch.listener = l
}

const enqueueTimeout = 5 * time.Minute

func (ch *channel[T]) Enqueue(item T) {
	ch.log("enqueue item")

	jj, err := json.Marshal(pgJob[T]{
		Channel: ch.name,
		Item:    item,
	})

	if err != nil {
		ch.log("error marshalling pgJob: %s", err.Error())
		return
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), enqueueTimeout)
	defer cancelCtx()

	conn, err := again.Retry[*pgxpool.Conn](ctx, func(ctx context.Context) (*pgxpool.Conn, error) {
		ch.log("trying to acquire connection")
		return ch.pool.Acquire(context.Background())
	})

	if err != nil {
		ch.log("error acquiring connection: %s", err.Error())
		return
	}
	ch.log("acquired connection for")
	defer conn.Release()

	_, err = conn.Query(ctx, fmt.Sprintf(`select pg_notify('%s', $1)`, ch.postgresQueueDriver.channelName), jj)
	if err != nil {
		ch.log("error notifying postgres: %s", err.Error())
		return
	}

	ch.log("notified postgres")
}

func (ch *channel[T]) Name() string {
	return ch.name
}
