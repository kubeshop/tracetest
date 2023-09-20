package testconnection

import (
	"fmt"
	"sync"
)

type NotifierFn func(Job)

type Listener struct {
	subscriptions map[string]NotifierFn
	mutex         sync.Mutex
}

func NewListener() *Listener {
	return &Listener{
		subscriptions: make(map[string]NotifierFn),
	}
}

func (m *Listener) Subscribe(jobID string, notifier NotifierFn) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.subscriptions[jobID] != nil {
		return fmt.Errorf("already subscribed to job %s", jobID)
	}

	m.subscriptions[jobID] = notifier
	return nil
}

func (m *Listener) Unsubscribe(jobID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.subscriptions, jobID)
}

func (m *Listener) Notify(job Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	notifierFn := m.subscriptions[job.ID]

	if notifierFn != nil {
		notifierFn(job)
	}
}
