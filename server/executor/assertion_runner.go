package executor

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/traces"
)

type AssertionRequest struct {
	TestDefinition assertions.TestDefinition
	TestRun        openapi.TestRun
}

type AssertionRunner interface {
	RunAssertions(request AssertionRequest)
	WorkerPool
}

type defaultAssertionRunner struct {
	db           testdb.Repository
	inputChannel chan AssertionRequest
	exitChannel  chan bool
}

var _ WorkerPool = &defaultAssertionRunner{}
var _ AssertionRunner = &defaultAssertionRunner{}

func NewAssertionRunner(db testdb.Repository) AssertionRunner {
	return &defaultAssertionRunner{
		db:           db,
		inputChannel: make(chan AssertionRequest, 1),
	}
}

func (e *defaultAssertionRunner) Start(workers int) {
	e.exitChannel = make(chan bool, workers)

	for i := 0; i < workers; i++ {
		ctx := context.Background()
		go e.startWorker(ctx)
	}
}

func (e *defaultAssertionRunner) Stop() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			e.exitChannel <- true
			return
		}
	}
}

func (e *defaultAssertionRunner) startWorker(ctx context.Context) {
	for {
		select {
		case <-e.exitChannel:
			fmt.Println("Exiting assertion executor worker")
			return
		case assertionRequest := <-e.inputChannel:
			err := e.runAssertionsAndUpdateResult(ctx, assertionRequest)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func (e *defaultAssertionRunner) runAssertionsAndUpdateResult(ctx context.Context, request AssertionRequest) error {
	response, err := e.executeAssertions(ctx, request)
	if err != nil {
		return err
	}

	err = e.db.UpdateRun(ctx, response)
	if err != nil {
		return fmt.Errorf("could not save result on database: %w", err)
	}

	return nil
}

func (e *defaultAssertionRunner) executeAssertions(ctx context.Context, request AssertionRequest) (*openapi.TestRun, error) {
	trace, err := traces.FromOtel(request.TestRun.Trace)
	if err != nil {
		return nil, err
	}

	test, err := e.db.GetTest(ctx, request.TestRun.TestId)
	if err != nil {
		return nil, err
	}

	// Temporary patch to disable the assertion engine if frontend request is not prepared yet (old selector format)
	if e.shouldIgnoreTest(test) {
		return &request.TestRun, nil
	}

	result := assertions.Assert(trace, request.TestDefinition)

	e.setResults(&request.TestRun, result)

	return &request.TestRun, nil
}

func (e *defaultAssertionRunner) shouldIgnoreTest(test *openapi.Test) bool {
	// If any assertion uses the old selector format, ignore the whole test and don't execute the
	// assertions.
	for selector, asserts := range test.Definition.Definitions {
		for _, assert := range asserts {
			if selector == "" && len(assert.Selectors) > 0 {
				return true
			}
		}
	}

	return false
}

func (e *defaultAssertionRunner) setResults(run *openapi.TestRun, testResult assertions.TestResult) {
	run.State = TestRunStateFinished
	run.CompletedAt = time.Now()

	res := openapi.AssertionResults{
		Results:   map[string][]openapi.AssertionResult{},
		AllPassed: true,
	}
	for selector, assertionResults := range testResult {
		res.Results[string(selector)] = make([]openapi.AssertionResult, len(assertionResults))
		for i, assertionResult := range assertionResults {
			res.Results[string(selector)][i] = openapi.AssertionResult{
				Id:          assertionResult.ID,
				Attribute:   assertionResult.Attribute,
				Comparator:  assertionResult.Comparator.String(),
				Expected:    assertionResult.Value,
				SpanResults: make([]openapi.AssertionSpanResult, len(assertionResult.AssertionSpanResults)),
			}
			for j, spanAssertionResult := range assertionResult.AssertionSpanResults {
				spanID := hex.EncodeToString(spanAssertionResult.Span.ID[:])
				testPassed := spanAssertionResult.CompareErr == nil
				if !testPassed {
					res.AllPassed = false
				}

				res.Results[string(selector)][i].SpanResults[j] = openapi.AssertionSpanResult{
					SpanId:        spanID,
					ObservedValue: spanAssertionResult.ActualValue,
					Passed:        testPassed,
				}
			}
		}
	}

	run.Result = res
}

func (e *defaultAssertionRunner) RunAssertions(request AssertionRequest) {
	e.inputChannel <- request
}
