package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

func NewPostgresQueueDriver(pool *pgxpool.Pool) *PostgresQueueDriver {
	id := id.GenerateID()
	return &PostgresQueueDriver{
		log:      newLoggerFn("PostgresQueueDriver - " + id.String()),
		pool:     pool,
		channels: map[string]*channel{},
		exit:     make(chan bool),
	}
}

// PostgresQueueDriver is a queue driver that uses Postgres LISTEN/NOTIFY
// Since each queue needs its own connection, it's not practical/scalable
// to create a new Driver instance for each queue. Instead, we create a
// single Driver instance and use it to create channels for each queue.
//
// This driver requires one connection that listens to messages in any queue
// and routes them to the correct worker.
type PostgresQueueDriver struct {
	log      loggerFn
	pool     *pgxpool.Pool
	channels map[string]*channel
	running  bool
	exit     chan bool
}

func (qd *PostgresQueueDriver) getChannel(name string) (*channel, error) {
	ch, ok := qd.channels[name]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", name)
	}

	return ch, nil
}

const pgChannelName = "tracetest_queue"

type pgJob struct {
	Channel string `json:"channel"`
	Job     Job    `json:"job"`
}

func (qd *PostgresQueueDriver) Start() {
	if qd.running {
		// we want only 1 worker here
		qd.log("already running")
		return
	}
	qd.running = true

	go func(qd *PostgresQueueDriver) {
		qd.log("start")

		qd.log("acquiring connection")
		conn, err := qd.pool.Acquire(context.Background())
		if err != nil {
			panic(fmt.Errorf("error acquiring connection: %w", err))
		}
		defer conn.Release()

		for {
			select {
			case <-qd.exit:
				qd.log("exit")
				return
			default:
				qd.worker(conn)
			}
		}
	}(qd)
}

func (qd *PostgresQueueDriver) worker(conn *pgxpool.Conn) {
	qd.log("listening for notifications")
	_, err := conn.Exec(context.Background(), "listen "+pgChannelName)
	if err != nil {
		qd.log("error listening for notifications: %s", err.Error())
		return
	}
	qd.log("waiting for notification")
	notification, err := conn.Conn().WaitForNotification(context.Background())
	if err != nil {
		qd.log("error waiting for notification: %s", err.Error())
		return
	}

	job := pgJob{}
	err = json.Unmarshal([]byte(notification.Payload), &job)
	if err != nil {
		qd.log("error unmarshalling pgJob: %s", err.Error())
		return
	}

	qd.log("received job for channel: %s, runID: %d", job.Channel, job.Job.Run.ID)

	channel, err := qd.getChannel(job.Channel)
	if err != nil {
		qd.log("error getting channel: %s", err.Error())
		return
	}

	// spin off so we can keep listening for jobs
	channel.listener.Listen(job.Job)
	qd.log("spun off job for channel: %s, runID: %d", job.Channel, job.Job.Run.ID)
}

func (qd *PostgresQueueDriver) Stop() {
	qd.log("stopping")
	qd.exit <- true
}

// Channel registers a new queue channel and returns it
func (qd *PostgresQueueDriver) Channel(name string) *channel {
	if _, channelNameExists := qd.channels[name]; channelNameExists {
		panic(fmt.Errorf("channel %s already exists", name))
	}

	ch := &channel{
		PostgresQueueDriver: qd,
		name:                name,
		log:                 newLoggerFn(fmt.Sprintf("PostgresQueueDriver - %s", name)),
		pool:                qd.pool,
	}

	qd.channels[name] = ch

	return ch
}

type channel struct {
	*PostgresQueueDriver
	name     string
	log      loggerFn
	pool     *pgxpool.Pool
	listener Listener
}

func (ch *channel) SetListener(l Listener) {
	ch.listener = l
}

const enqueueTimeout = 5 * time.Second

func (ch *channel) Enqueue(job Job) {
	ch.log("enqueue job for run %d", job.Run.ID)

	jj, err := json.Marshal(pgJob{
		Channel: ch.name,
		Job:     job,
	})

	if err != nil {
		ch.log("error marshalling pgJob: %s", err.Error())
		return
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), enqueueTimeout)
	defer cancelCtx()

	conn, err := ch.pool.Acquire(context.Background())
	if err != nil {
		ch.log("error acquiring connection: %s", err.Error())
		return
	}
	ch.log("aquired connection for run %d", job.Run.ID)
	defer conn.Release()

	_, err = conn.Query(ctx, fmt.Sprintf(`select pg_notify('%s', $1)`, pgChannelName), jj)
	if err != nil {
		ch.log("error notifying postgres: %s", err.Error())
		return
	}

	ch.log("notified postgres for run %d", job.Run.ID)
}
