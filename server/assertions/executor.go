package assertions

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/traces"
)

type RunAssertionsMessage struct {
	Test   openapi.Test
	Result openapi.TestRunResult
}

type Executor struct {
	inputChannel, outputChannel chan RunAssertionsMessage
}

func NewExecutor(inputChannel, outputChannel chan RunAssertionsMessage) Executor {
	return Executor{
		inputChannel:  inputChannel,
		outputChannel: outputChannel,
	}
}

func (e Executor) Start() {
	for {
		select {
		case request := <-e.inputChannel:
			response, err := e.executeAssertions(request)
			if err != nil {
				fmt.Println(err)
			}

			e.outputChannel <- *response
		}
	}
}

func (e Executor) executeAssertions(request RunAssertionsMessage) (*RunAssertionsMessage, error) {
	trace, err := traces.FromOtel(request.Result.Trace)
	if err != nil {
		return nil, err
	}

	testDefinition := convertAssertionsIntoTestDefinition(request.Test.Assertions)

	result := Assert(trace, testDefinition)

	response := e.setResults(request, result)

	return response, nil
}

func (e Executor) setResults(request RunAssertionsMessage, testResult TestResult) *RunAssertionsMessage {
	response := request
	response.Result.State = executor.TestRunStateFinished
	response.Result.CompletedAt = time.Now()
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

	response.Result.AssertionResult = assertionResultArray
	response.Result.AssertionResultState = allTestsPassed

	return &response
}
