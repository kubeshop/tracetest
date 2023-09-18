package testconnection

import (
	"sync"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
)

type Job struct {
	ID         string
	DataStore  datastore.DataStore
	TestResult model.ConnectionResult
}

type NotifierFn func(Job)

type Listener struct {
	subscriptions map[string][]NotifierFn
	mutex         sync.Mutex
}

func NewListener() *Listener {
	return &Listener{
		subscriptions: make(map[string][]NotifierFn),
	}
}

func (m *Listener) Subscribe(jobID string, notifier NotifierFn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.subscriptions[jobID] = append(m.subscriptions[jobID], notifier)
}

func (m *Listener) Unsubscribe(jobID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.subscriptions, jobID)
}

func (m *Listener) Notify(job Job) {
	for _, sub := range m.subscriptions[job.ID] {
		sub(job)
	}
}
