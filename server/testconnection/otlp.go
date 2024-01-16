package testconnection

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/subscription"
)

type OTLPConnectionTester struct {
	subscriptionManager   subscription.Manager
	connectionInformation map[string]connectionInformation
	subscriptions         map[string]subscription.Subscription
}

type OTLPConnectionTestRequest struct{}

type OTLPConnectionTestResponse struct {
	TenantID string

	NumberTraces int
}

type connectionInformation struct {
	ReceivedTraces bool
}

func NewOTLPConnectionTester(subscriptionManager subscription.Manager) *OTLPConnectionTester {
	return &OTLPConnectionTester{subscriptionManager: subscriptionManager}
}

func (t *OTLPConnectionTester) StartTest(ctx context.Context) {
	tenantID := middleware.TenantIDFromContext(ctx)
	subject := fmt.Sprintf("start_otlp_connection_test_%s", tenantID)
	responseSubject := fmt.Sprintf("otlp_connection_test_incoming_spans_%s", tenantID)
	t.subscriptionManager.Publish(subject, OTLPConnectionTestRequest{})

	t.subscriptionManager.Subscribe(responseSubject, subscription.NewSubscriberFunction(func(m subscription.Message) error {

	}))
}

func (t *OTLPConnectionTester) EndTest(ctx context.Context) {
	tenantID := middleware.TenantIDFromContext(ctx)
	subject := fmt.Sprintf("end_otlp_connection_test_%s", tenantID)
	t.subscriptionManager.Publish(subject, OTLPConnectionTestRequest{})
}
