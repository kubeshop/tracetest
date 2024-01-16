package testconnection_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testconnection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccessfulGetSpanCount(t *testing.T) {
	manager := subscription.NewManager()
	tester := testconnection.NewOTLPConnectionTester(manager)

	go func() {
		time.Sleep(time.Second)
		manager.Publish(testconnection.PostSpanCountTopicName(), testconnection.OTLPConnectionTestResponse{
			NumberSpans:       2,
			LastSpanTimestamp: time.Now(),
		})
	}()

	response, err := tester.GetSpanCount(context.Background(), testconnection.WithTimeout(10*time.Second))
	require.NoError(t, err)

	assert.Equal(t, 2, response.NumberSpans)
	assert.False(t, response.LastSpanTimestamp.IsZero())
}

func TestGetSpanCountTimeout(t *testing.T) {
	// Given that GetSpanCount is called but the "PostSpanCountTopicName" topic never receives a message
	// the GetSpanCount call should timeout

	manager := subscription.NewManager()
	tester := testconnection.NewOTLPConnectionTester(manager)

	_, err := tester.GetSpanCount(context.Background(), testconnection.WithTimeout(1*time.Second))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestSuccessfulGetSpanCountWithTenantID(t *testing.T) {
	tenantID := uuid.NewString()
	ctx := context.WithValue(context.Background(), middleware.TenantIDKey, tenantID)
	manager := subscription.NewManager()
	tester := testconnection.NewOTLPConnectionTester(manager)

	go func() {
		time.Sleep(time.Second)
		manager.Publish(testconnection.PostSpanCountTopicName(testconnection.WithTenantID(tenantID)), testconnection.OTLPConnectionTestResponse{
			TenantID:          tenantID,
			NumberSpans:       2,
			LastSpanTimestamp: time.Now(),
		})
	}()

	response, err := tester.GetSpanCount(ctx, testconnection.WithTimeout(10*time.Second))
	require.NoError(t, err)

	assert.Equal(t, tenantID, response.TenantID)
	assert.Equal(t, 2, response.NumberSpans)
	assert.False(t, response.LastSpanTimestamp.IsZero())
}
