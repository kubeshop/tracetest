package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/kubeshop/tracetest/agent/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	ReconnectRetryAttempts     = 6
	ReconnectRetryAttemptDelay = 1 * time.Second
	defaultPingPeriod          = 30 * time.Second
)

type Config struct {
	APIKey     string
	AgentName  string
	PingPeriod time.Duration
}

type SessionConfig struct {
	BatchTimeout        time.Duration
	AgentIdentification *proto.AgentIdentification
}

type Client struct {
	mutex         sync.Mutex
	endpoint      string
	conn          *grpc.ClientConn
	config        Config
	sessionConfig *SessionConfig
	done          chan bool

	logger *zap.Logger

	triggerListener             func(context.Context, *proto.TriggerRequest) error
	pollListener                func(context.Context, *proto.PollingRequest) error
	shutdownListener            func(context.Context, *proto.ShutdownRequest) error
	dataStoreConnectionListener func(context.Context, *proto.DataStoreConnectionTestRequest) error
}

func (c *Client) Start(ctx context.Context) error {
	err := c.startup(ctx)
	if err != nil {
		return err
	}

	c.done = make(chan bool)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c.done
		// We cannot `defer cancel()` in this case because the start listener functions
		// start one goroutine each and don't block the execution of this function.
		// Thus, if we cancel the context, all those goroutines will fail.
		cancel()
	}()

	err = c.startTriggerListener(ctx)
	if err != nil {
		return err
	}

	err = c.startPollerListener(ctx)
	if err != nil {
		return err
	}

	err = c.startShutdownListener(ctx)
	if err != nil {
		return err
	}

	err = c.startDataStoreConnectionTestListener(ctx)
	if err != nil {
		return err
	}

	err = c.startHearthBeat(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) WaitUntilDisconnected() {
	<-c.done
}

func (c *Client) SessionConfiguration() *SessionConfig {
	if c.sessionConfig == nil {
		return nil
	}

	deferredPtr := *c.sessionConfig
	return &deferredPtr
}

func (c *Client) Close() error {
	c.done <- true
	return c.conn.Close()
}

func (c *Client) OnTriggerRequest(listener func(context.Context, *proto.TriggerRequest) error) {
	c.triggerListener = listener
}

func (c *Client) OnDataStoreTestConnectionRequest(listener func(context.Context, *proto.DataStoreConnectionTestRequest) error) {
	c.dataStoreConnectionListener = listener
}

func (c *Client) OnPollingRequest(listener func(context.Context, *proto.PollingRequest) error) {
	c.pollListener = listener
}

func (c *Client) OnConnectionClosed(listener func(context.Context, *proto.ShutdownRequest) error) {
	c.shutdownListener = listener
}

func (c *Client) getConnectionRequest() (*proto.ConnectRequest, error) {
	name, err := c.getName()
	if err != nil {
		return nil, err
	}

	request := proto.ConnectRequest{
		ApiKey: c.config.APIKey,
		Name:   name,
	}

	return &request, nil
}

// getName retrieves the name of the agent. By default, it is the host name, however,
// it can be overwritten with an environment variable, or a flag.
func (c *Client) getName() (string, error) {
	if name := c.config.AgentName; name != "" {
		return name, nil
	}

	if name := os.Getenv("TRACETEST_AGENT_NAME"); name != "" {
		return name, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get hostname: %w", err)
	}

	return hostname, nil
}

func isCancelledError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "context canceled")
}

func (c *Client) reconnect() error {
	if c.conn != nil {
		c.conn.Close()
	}

	// connection is not working. We need to reconnect
	err := retry.Do(func() error {
		return c.connect(context.Background())
	}, retry.Attempts(ReconnectRetryAttempts), retry.Delay(ReconnectRetryAttemptDelay))

	if err != nil {
		return fmt.Errorf("could not reconnect to server: %w", err)
	}

	return c.Start(context.Background())
}

func (c *Client) handleDisconnectionError(inputErr error) (bool, error) {
	if !isConnectionError(inputErr) {
		// if it's nil or any error other than the one we care about, return it and let the caller handle it
		return false, inputErr
	}

	log.Printf("Reconnecting agent due to error: %s", inputErr.Error())
	err := retry.Do(func() error {
		return c.reconnect()
	})

	if err != nil {
		log.Fatal(err)
	}

	return true, nil
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	possibleErrors := []string{
		"connection refused",
		"server closed",
		"token is expired",

		// From time to time, the server can start sending those errors to the
		// agent. This mitigates the risk of an agent getting stuck in an error state
		"unexpected HTTP status code received from server: 500",
		"the client connection is closing",
	}
	for _, possibleErr := range possibleErrors {
		if strings.Contains(err.Error(), possibleErr) {
			return true
		}
	}

	return false
}

func isEndOfFileError(err error) bool {
	if err == nil {
		return false
	}

	if isEof := errors.Is(err, io.EOF); isEof {
		return true
	}

	return strings.Contains(err.Error(), "EOF")
}
