package openapi

type WorkerPool interface {
	Start(workers int)
	Stop()
}
