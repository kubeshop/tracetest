package utils

func RunStateIsFinished(state string) bool {
	return RunStateIsFailed(state) || state == "FINISHED"
}

func RunStateIsFailed(state string) bool {
	return state == "TRIGGER_FAILED" ||
		state == "TRACE_FAILED" ||
		state == "ASSERTION_FAILED" ||
		state == "ANALYZING_ERROR" ||
		state == "FAILED" // this one is for backwards compatibility
}
