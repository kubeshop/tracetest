package testconnection

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/subscription"
	"golang.org/x/sync/semaphore"
)

type TopicNameOption func(*topicNameConfig)

type topicNameConfig struct {
	TenantID string
}

func WithTenantID(tenantID string) TopicNameOption {
	return func(tnc *topicNameConfig) {
		tnc.TenantID = tenantID
	}
}

func GetSpanCountTopicName(opts ...TopicNameOption) string {
	var config topicNameConfig
	for _, opt := range opts {
		opt(&config)
	}

	return fmt.Sprintf("otlp_connection_test_get_span_count_%s", config.TenantID)
}
func PostSpanCountTopicName(opts ...TopicNameOption) string {
	var config topicNameConfig
	for _, opt := range opts {
		opt(&config)
	}

	return fmt.Sprintf("otlp_connection_test_span_count_%s", config.TenantID)
}
func ResetSpanCountTopicName(opts ...TopicNameOption) string {
	var config topicNameConfig
	for _, opt := range opts {
		opt(&config)
	}

	return fmt.Sprintf("otlp_connection_test_reset_span_count_%s", config.TenantID)
}

type OTLPConnectionTester struct {
	subscriptionManager subscription.Manager
}

type OTLPConnectionTestRequest struct{}

type OTLPConnectionTestResponse struct {
	TenantID string

	NumberSpans       int
	LastSpanTimestamp time.Time
}

func NewOTLPConnectionTester(subscriptionManager subscription.Manager) *OTLPConnectionTester {
	return &OTLPConnectionTester{subscriptionManager: subscriptionManager}
}

type GetSpanCountOption func(*getSpanCountConfig)
type getSpanCountConfig struct {
	timeout time.Duration
}

func WithTimeout(timeout time.Duration) GetSpanCountOption {
	return func(gscc *getSpanCountConfig) {
		gscc.timeout = timeout
	}
}

func (t *OTLPConnectionTester) GetSpanCount(ctx context.Context, opts ...GetSpanCountOption) (OTLPConnectionTestResponse, error) {
	config := getSpanCountConfig{
		timeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(&config)
	}

	ctx, cancel := context.WithTimeout(ctx, config.timeout)
	defer cancel()

	semaphore := semaphore.NewWeighted(1)
	tenantID := middleware.TenantIDFromContext(ctx)
	t.subscriptionManager.Publish(GetSpanCountTopicName(WithTenantID(tenantID)), OTLPConnectionTestRequest{})

	var response OTLPConnectionTestResponse
	semaphore.Acquire(ctx, 1)
	topicName := PostSpanCountTopicName(WithTenantID(tenantID))
	subscriber := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		m.DecodeContent(&response)
		semaphore.Release(1)
		return nil
	})

	t.subscriptionManager.Subscribe(topicName, subscriber)
	defer t.subscriptionManager.Unsubscribe(topicName, subscriber.ID())

	// Acts as a WaitGroup.Wait() that is canceled with the context in case of timeout.
	err := semaphore.Acquire(ctx, 1)
	if err != nil {
		return OTLPConnectionTestResponse{}, fmt.Errorf("could not get span count: %w", err)
	}

	return response, nil
}
