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

type AssertionFinishCallback func(openapi.Test, openapi.TestRunResult)

type AssertionRequest struct {
	TestDefinition assertions.TestDefinition
	Result         openapi.TestRunResult
}

type AssertionRunner interface {
	RunAssertions(ctx context.Context, result openapi.TestRunResult) error
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

	err = e.db.UpdateResult(ctx, response)
	if err != nil {
		return fmt.Errorf("could not save result on database: %w", err)
	}

	return nil
}

func (e *defaultAssertionRunner) executeAssertions(ctx context.Context, request AssertionRequest) (*openapi.TestRunResult, error) {
	trace, err := traces.FromOtel(request.Result.Trace)
	if err != nil {
		return nil, err
	}

	test, err := e.db.GetTest(ctx, request.Result.TestId)
	if err != nil {
		return nil, err
	}

	// Temporary patch to disable the assertion engine if frontend request is not prepared yet (old selector format)
	if e.shouldIgnoreTest(test) {
		return &request.Result, nil
	}

	result := assertions.Assert(trace, request.TestDefinition)

	e.setResults(&request.Result, result)

	return &request.Result, nil
}

func (e *defaultAssertionRunner) shouldIgnoreTest(test *openapi.Test) bool {
	// If any assertion uses the old selector format, ignore the whole test and don't execute the
	// assertions.
	for _, assertion := range test.Assertions {
		if assertion.Selector == "" && len(assertion.Selectors) > 0 {
			return true
		}
	}

	return false
}

func (e *defaultAssertionRunner) setResults(result *openapi.TestRunResult, testResult assertions.TestResult) {
	result.State = TestRunStateFinished
	result.CompletedAt = time.Now()
	assertionResultArray := make([]openapi.AssertionResult, 0)
	allTestsPassed := true

	for _, assertionExecutionResult := range testResult {
		for _, assertionResult := range assertionExecutionResult {
			spanAssertions := make([]openapi.SpanAssertionResult, 0)
			for _, spanAssertionResult := range assertionResult.AssertionSpanResults {
				spanID := hex.EncodeToString(spanAssertionResult.Span.ID[:])
				testPassed := spanAssertionResult.CompareErr == nil
				if !testPassed {
					allTestsPassed = false
				}

				spanAssertions = append(spanAssertions, openapi.SpanAssertionResult{
					Passed:        testPassed,
					SpanId:        spanID,
					ObservedValue: spanAssertionResult.ActualValue,
				})
			}

			result := openapi.AssertionResult{
				AssertionId:          assertionResult.Assertion.ID,
				SpanAssertionResults: spanAssertions,
			}

			assertionResultArray = append(assertionResultArray, result)
		}
	}

	result.AssertionResult = assertionResultArray
	result.AssertionResultState = allTestsPassed
}

func (e *defaultAssertionRunner) RunAssertions(ctx context.Context, result openapi.TestRunResult) error {
	test, err := e.db.GetTest(ctx, result.TestId)
	if err != nil {
		return err
	}

	testDefinition, err := ConvertAssertionsIntoTestDefinition(test.Assertions)
	if err != nil {
		return err
	}

	assertionRequest := AssertionRequest{
		TestDefinition: testDefinition,
		Result:         result,
	}

	e.inputChannel <- assertionRequest

	return nil
}
