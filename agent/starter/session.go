package starter

import "github.com/kubeshop/tracetest/agent/client"

type Session struct {
	Token  string
	client *client.Client
}

func (s *Session) Close() {
	s.client.Close()
}

func (s *Session) WaitUntilDisconnected() {
	s.client.WaitUntilDisconnected()
}
