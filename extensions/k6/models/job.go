package models

import (
	"fmt"

	"github.com/kubeshop/tracetest/extensions/k6/openapi"
)

type JobType string

const (
	RunTestFromId JobType = "runTestFromId"
)

type JobStatus string

const (
	Pending JobStatus = "pending"
	Running JobStatus = "running"
	Failed  JobStatus = "failed"
	Success JobStatus = "success"
)

type Job struct {
	TraceID    string
	TestID     string
	JobType    JobType
	Request    Request
	Run        *TracetestRun
	JobStatus  JobStatus
	ShouldWait bool
}

func NewJob(traceID, testID string, jobType JobType, shouldWait bool, request Request) Job {
	return Job{
		TraceID:    traceID,
		TestID:     testID,
		JobType:    jobType,
		Request:    request,
		JobStatus:  Pending,
		ShouldWait: shouldWait,
	}
}

func (job Job) HandleRunResponse(run *openapi.TestRun, err error) Job {
	if run == nil {
		job.JobStatus = Failed
	} else {
		job.JobStatus = Success
		job.Run = &TracetestRun{
			TestRun: run,
			TestId:  job.TestID,
		}
	}

	return job
}

func (job Job) Summary(baseUrl string) string {
	runSummary := "JobStatus=" + string(job.JobStatus)
	if job.Run != nil {
		runSummary = job.Run.Summary(baseUrl)
	}

	return fmt.Sprintf("Request=%s - %s, TraceID=%s, %s", job.Request.Method, job.Request.URL, job.TraceID, runSummary)
}

func (job Job) IsSuccessful() bool {
	isJobStatusSuccessful := job.JobStatus == Success
	runExists := job.Run != nil

	return isJobStatusSuccessful && runExists && job.Run.IsSuccessful()
}
