package tracetest

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/extensions/k6/models"
	"github.com/kubeshop/tracetest/extensions/k6/openapi"
)

func NewAPIClient(options Options) *openapi.APIClient {
	url, err := url.Parse(options.ServerUrl)

	if err != nil {
		panic(err)
	}

	config := openapi.NewConfiguration()
	config.Host = url.Host
	config.Scheme = url.Scheme

	if options.ServerPath != "" {
		config.Servers = []openapi.ServerConfiguration{
			{
				URL: options.ServerPath,
			},
		}
	}

	return openapi.NewAPIClient(config)
}

func (t *Tracetest) runTest(testID, traceId string) (*openapi.TestRun, error) {
	request := t.client.ApiApi.RunTest(context.Background(), testID)
	key := "TRACE_ID"
	request = request.RunInformation(openapi.RunInformation{
		Variables: []openapi.EnvironmentValue{{
			Key:   &key,
			Value: &traceId,
		}},
	})

	run, _, err := t.client.ApiApi.RunTestExecute(request)
	return run, err
}

func (t *Tracetest) waitForTestResult(testID, testRunID string) (openapi.TestRun, error) {
	var (
		testRun   openapi.TestRun
		lastError error
		wg        sync.WaitGroup
	)

	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				readyTestRun, err := t.getIsTestReady(context.Background(), testID, testRunID)
				if err != nil {
					lastError = err
					wg.Done()
					return
				}

				if readyTestRun != nil {
					testRun = *readyTestRun
					wg.Done()
					return
				}
			}
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.TestRun{}, lastError
	}

	return testRun, nil
}

func (t *Tracetest) getIsTestReady(ctx context.Context, testID, testRunId string) (*openapi.TestRun, error) {
	req := t.client.ApiApi.GetTestRun(ctx, testID, testRunId)
	run, _, err := t.client.ApiApi.GetTestRunExecute(req)

	if err != nil {
		return &openapi.TestRun{}, fmt.Errorf("could not execute GetTestRun request: %w", err)
	}

	if *run.State == "FAILED" || *run.State == "FINISHED" {
		return run, nil
	}

	return nil, nil
}

func (t *Tracetest) stringSummary() string {
	failedSummary := "[FAILED] \n"
	successfulSummary := "[SUCCESSFUL] \n"
	totalRuns := 0
	failedRuns := 0
	successfulRuns := 0

	t.processedBuffer.Range(func(_, value interface{}) bool {
		if job, ok := value.(models.Job); ok {
			totalRuns += 1
			if job.IsSuccessful() {
				successfulSummary += fmt.Sprintf("[%s] \n", job.Summary(t.options.ServerUrl))
				successfulRuns += 1
			} else {
				failedSummary += fmt.Sprintf("[%s] \n", job.Summary(t.options.ServerUrl))
				failedRuns += 1
			}
		}

		return true
	})

	totalResults := fmt.Sprintf("[TotalRuns=%d, SuccessfulRus=%d, FailedRuns=%d] \n", totalRuns, successfulRuns, failedRuns)

	if failedRuns == 0 {
		failedSummary = ""
	}

	if successfulRuns == 0 {
		successfulSummary = ""
	}

	return totalResults + failedSummary + successfulSummary
}

type JsonResult struct {
	TotalRuns      int
	SuccessfulRuns int
	FailedRuns     int
	Failed         []models.Job
	Successful     []models.Job
}

func (t *Tracetest) jsonSummary() JsonResult {
	JsonResult := JsonResult{
		TotalRuns:      0,
		SuccessfulRuns: 0,
		FailedRuns:     0,
		Failed:         []models.Job{},
		Successful:     []models.Job{},
	}

	t.processedBuffer.Range(func(_, value interface{}) bool {
		if job, ok := value.(models.Job); ok {
			JsonResult.TotalRuns += 1
			if job.IsSuccessful() {
				JsonResult.Successful = append(JsonResult.Successful, job)
				JsonResult.SuccessfulRuns += 1
			} else {
				JsonResult.Failed = append(JsonResult.Failed, job)
				JsonResult.FailedRuns += 1
			}
		}

		return true
	})

	return JsonResult
}
