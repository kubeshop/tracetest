package executor

type WorkerPool interface {
	Start(workers int)
	Stop()
}
