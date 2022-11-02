package model

import "encoding/json"

func BumpTransactionVersionIfNeeded(in, updated Transaction) Transaction {
	transactionHasChanged := transactionHasChanged(in, updated)
	if transactionHasChanged {
		updated.Version = in.Version + 1
	}

	return updated
}

func transactionHasChanged(in, updated Transaction) bool {
	jsons := []struct {
		Name        string
		Description string
		Steps       []string
	}{
		{
			Name:        in.Name,
			Description: in.Description,
			Steps:       getStepIds(in),
		},
		{
			Name:        updated.Name,
			Description: updated.Description,
			Steps:       getStepIds(updated),
		},
	}

	inJson, _ := json.Marshal(jsons[0])
	updatedJson, _ := json.Marshal(jsons[1])

	return string(inJson) != string(updatedJson)
}

func getStepIds(in Transaction) []string {
	steps := make([]string, len(in.Steps))
	for i, step := range in.Steps {
		steps[i] = step.ID.String()
	}

	return steps
}
