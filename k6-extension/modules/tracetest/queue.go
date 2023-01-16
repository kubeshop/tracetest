package tracetest

type JobType string

const (
	RunTestFromId         JobType = "runTestFromId"
	RunTestFromDefinition JobType = "runTestFromDefinition"
)

type Metadata map[string]string

type Job struct {
	traceID    string
	testID     string
	definition string
	jobType    JobType
}

func (t *Tracetest) processQueue() {
	t.bufferLock.Lock()
	bufferedJobs := t.buffer
	t.buffer = make([]Job, 0, len(bufferedJobs)) // flushing queue
	t.bufferLock.Unlock()

	for _, job := range bufferedJobs {
		switch job.jobType {
		case RunTestFromId:
			t.runFromId(job.testID, job.traceID)
		case RunTestFromDefinition:
			t.runFromDefinition(job.definition, job.traceID)
		}
	}
}

func (t *Tracetest) queueJob(job Job) {
	t.logger.Infoln("Queuing Job: ", job)
	t.bufferLock.Lock()
	defer t.bufferLock.Unlock()

	t.buffer = append(t.buffer, job)
}
