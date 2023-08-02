package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type loggerFn func(string, ...any)

func newLoggerFn(name string) loggerFn {
	return func(format string, params ...any) {
		log.Printf("[%s] %s", name, fmt.Sprintf(format, params...))
	}
}

func NewInMemoryQueueDriver(name string) *InMemoryQueueDriver {
	return &InMemoryQueueDriver{
		log:   newLoggerFn(fmt.Sprintf("InMemoryQueueDriver - %s", name)),
		queue: make(chan Job),
		exit:  make(chan bool),
		name:  name,
	}
}

type InMemoryQueueDriver struct {
	log   loggerFn
	queue chan Job
	exit  chan bool
	q     *Queue
	name  string
}

func (qd *InMemoryQueueDriver) SetQueue(q *Queue) {
	qd.q = q
}

func (qd InMemoryQueueDriver) Enqueue(job Job) {
	qd.queue <- job
}

const inMemoryQueueWorkerCount = 5

func (qd InMemoryQueueDriver) Start() {
	for i := 0; i < inMemoryQueueWorkerCount; i++ {
		go func() {
			qd.log("start")
			for {
				select {
				case <-qd.exit:
					qd.log("exit")
					return
				case job := <-qd.queue:
					qd.q.Listen(job)
				}
			}
		}()
	}
}

func (qd InMemoryQueueDriver) Stop() {
	qd.log("stopping")
	qd.exit <- true
}

func NewPostgresQueueDriver(pool *pgxpool.Pool, name string) *PostgresQueueDriver {
	return &PostgresQueueDriver{
		log:  newLoggerFn(fmt.Sprintf("PostgresQueueDriver - %s", name)),
		pool: pool,
		exit: make(chan bool),
		name: name,
	}
}

type PostgresQueueDriver struct {
	log  loggerFn
	pool *pgxpool.Pool
	name string
	q    *Queue
	exit chan bool
}

func (qd *PostgresQueueDriver) SetQueue(q *Queue) {
	qd.q = q
}

func (qd *PostgresQueueDriver) Enqueue(job Job) {
	qd.log("enqueue")
	jj, err := json.Marshal(job)
	if err != nil {
		qd.log("error marshalling job: %s", err.Error())
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelCtx()

	_, err = qd.pool.Query(ctx, fmt.Sprintf(`select pg_notify('%s', $1)`, qd.name), jj)
	if err != nil {
		qd.log("error notifying postgres: %s", err.Error())
		return
	}

	qd.log("notified postgres")
}

func (qd *PostgresQueueDriver) Start() {
	go func(qd *PostgresQueueDriver) {
		qd.log("start")
		conn, err := qd.pool.Acquire(context.Background())
		if err != nil {
			panic(fmt.Errorf("error acquiring connection: %w", err))
		}
		defer conn.Release()

		_, err = conn.Exec(context.Background(), "listen "+qd.name)
		if err != nil {
			panic(fmt.Errorf("error listening: %w", err))
		}

		for {
			select {
			case <-qd.exit:
				qd.log("exit")
				return
			default:
				notification, err := conn.Conn().WaitForNotification(context.Background())
				if err != nil {
					qd.log("error waiting for notification: %s", err.Error())
				}

				job := Job{}
				err = json.Unmarshal([]byte(notification.Payload), &job)
				if err != nil {
					qd.log("error unmarshalling job: %s", err.Error())
				}

				qd.q.Listen(job)
			}
		}
	}(qd)
}

func (qd *PostgresQueueDriver) Stop() {
	qd.log("stopping")
	qd.exit <- true
}
