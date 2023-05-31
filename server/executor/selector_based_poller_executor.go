package executor

import "github.com/kubeshop/tracetest/server/model"

type SelectorBasedPollerExecutor struct {
	pollerExecutor PollerExecutor
}

func (pe SelectorBasedPollerExecutor) ExecuteRequest(request *PollingRequest) (bool, string, model.Run, error) {
	ready, reason, run, err := pe.pollerExecutor.ExecuteRequest(request)
	// if !ready {
	return ready, reason, run, err
	// }
}
