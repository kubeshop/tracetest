package executor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestPersistentPoller(t *testing.T) {
	t.Run("CanTracePolls", func(t *testing.T) {
		// NOT IMPLEMENTED YET
		t.Skip()
		res := openapi.TestRunResult{
			TestId:  "test",
			TraceId: "trace",
		}

		f := tracePollerSetup(t)

		f.mockTraceFetcher.expectTracePoll(res.TraceId)

		f.run([]openapi.TestRunResult{res}, 10*time.Millisecond)

		result := f.mockResultsDB.results[res.TestId]
		require.NotNil(t, result)
		assert.Equal(t, sampleParsedTraces, result.Trace)

		f.assert(t)
	})

}

var (
	sampleParsedTraces = openapi.ApiV3SpansResponseChunk{}
	sampleRawTraces    = &v1.TracesData{}
)

type tracePollerFixture struct {
	tracePoller      executor.PersistentTracePoller
	mockTraceFetcher *mockTraceFetcher
	mockResultsDB    *mockResultsDB
}

func (f tracePollerFixture) run(res []openapi.TestRunResult, ttl time.Duration) {
	f.tracePoller.Start(2)
	time.Sleep(10 * time.Millisecond)
	for _, r := range res {
		f.tracePoller.Poll(context.TODO(), r)
	}
	time.Sleep(ttl)
	f.tracePoller.Stop()
}

func (f tracePollerFixture) assert(t *testing.T) {
	f.mockTraceFetcher.AssertExpectations(t)
	f.mockResultsDB.AssertExpectations(t)
}

func tracePollerSetup(t *testing.T) tracePollerFixture {
	mtf := new(mockTraceFetcher)
	mtf.t = t
	mtf.Test(t)

	mr := new(mockResultsDB)
	mr.t = t
	mr.Test(t)

	mt := new(mockTestDB)
	mt.t = t
	mt.Test(t)

	mar := new(mockAssertionRunner)
	mar.t = t
	mar.Test(t)

	return tracePollerFixture{
		tracePoller:      executor.NewTracePoller(mtf, mr, mt, 100*time.Millisecond, subscription.NewManager(), mar),
		mockTraceFetcher: mtf,
		mockResultsDB:    mr,
	}
}

type mockTraceFetcher struct {
	mock.Mock
	t *testing.T
}

func (m *mockTraceFetcher) GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error) {
	args := m.Called(traceID)
	return args.Get(0).(*v1.TracesData), args.Error(1)
}

func (m *mockTraceFetcher) expectTracePoll(traceID string) *mock.Call {
	return m.
		On("GetTraceByID", traceID).
		Return(sampleRawTraces, noError)
}
