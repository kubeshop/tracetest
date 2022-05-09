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

type AssertionRunner interface {
	RunAssertions(result openapi.TestRunResult)
	WorkerPool
}

type defaultAssertionRunner struct {
	testDB       testdb.TestRepository
	resultDB     testdb.ResultRepository
	inputChannel chan openapi.TestRunResult
	exitChannel  chan bool
}

var _ WorkerPool = &defaultAssertionRunner{}
var _ AssertionRunner = &defaultAssertionRunner{}

func NewAssertionRunner(
	testRepository testdb.TestRepository,
	resultRepository testdb.ResultRepository,
	inputChannel chan openapi.TestRunResult,
) AssertionRunner {
	return &defaultAssertionRunner{
		testDB:       testRepository,
		resultDB:     resultRepository,
		inputChannel: inputChannel,
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
		case testResult := <-e.inputChannel:
			response, err := e.executeAssertions(ctx, testResult)
			if err != nil {
				fmt.Println(err.Error())
			}

			err = e.resultDB.UpdateResult(ctx, response)
			if err != nil {
				fmt.Println(fmt.Errorf("could not save result on database: %w", err).Error())
			}
		}
	}
}

func (e *defaultAssertionRunner) executeAssertions(ctx context.Context, testResult openapi.TestRunResult) (*openapi.TestRunResult, error) {
	trace, err := traces.FromOtel(testResult.Trace)
	if err != nil {
		return nil, err
	}

	test, err := e.testDB.GetTest(ctx, testResult.TestId)
	if err != nil {
		return nil, err
	}

	testDefinition := convertAssertionsIntoTestDefinition(test.Assertions)

	result := assertions.Assert(trace, testDefinition)

	e.setResults(&testResult, result)

	return &testResult, nil
}

func (e *defaultAssertionRunner) setResults(result *openapi.TestRunResult, testResult assertions.TestResult) {
	result.State = TestRunStateFinished
	result.CompletedAt = time.Now()
	assertionResultArray := make([]openapi.AssertionResult, 0)
	allTestsPassed := true

	for _, assertionResult := range testResult {
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

	result.AssertionResult = assertionResultArray
	result.AssertionResultState = allTestsPassed
}

func (e *defaultAssertionRunner) RunAssertions(result openapi.TestRunResult) {
	e.inputChannel <- result
}
