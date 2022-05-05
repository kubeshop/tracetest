package assertions

import (
	"fmt"

	"github.com/kubeshop/tracetest/openapi"
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
			err := e.executeAssertions(request)
			if err != nil {
				fmt.Println(err)
			}

			e.outputChannel <- request
		}
	}
}

func (e Executor) executeAssertions(request RunAssertionsMessage) error {
	trace, err := convertOTelTraceIntoTraceTree(request.Result.Trace)
	if err != nil {
		return err
	}

	testDefinition := convertAssertionsIntoTestDefinition(request.Test.Assertions)

	result := Assert(trace, testDefinition)

	fmt.Println(result)

	return nil
}
