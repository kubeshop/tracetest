package model

import (
	"encoding/json"
)

func BumpTestVersionIfNeeded(in, updated Test) (Test, error) {
	testHasChanged, err := testHasChanged(in, updated)
	if err != nil {
		return Test{}, err
	}

	if testHasChanged {
		updated.Version = in.Version + 1
	}

	return updated, nil
}

func testHasChanged(oldTest Test, newTest Test) (bool, error) {
	outputsHaveChanged, err := testFieldHasChanged(oldTest.Outputs, newTest.Outputs)
	if err != nil {
		return false, err
	}

	definitionHasChanged, err := testFieldHasChanged(oldTest.Specs, newTest.Specs)
	if err != nil {
		return false, err
	}

	serviceUnderTestHasChanged, err := testFieldHasChanged(oldTest.ServiceUnderTest, newTest.ServiceUnderTest)
	if err != nil {
		return false, err
	}

	nameHasChanged := oldTest.Name != newTest.Name
	descriptionHasChanged := oldTest.Description != newTest.Description

	return outputsHaveChanged || definitionHasChanged || serviceUnderTestHasChanged || nameHasChanged || descriptionHasChanged, nil
}

func testFieldHasChanged(oldField interface{}, newField interface{}) (bool, error) {
	oldFieldJSON, err := json.Marshal(oldField)
	if err != nil {
		return false, err
	}

	newFieldJSON, err := json.Marshal(newField)
	if err != nil {
		return false, err
	}

	return string(oldFieldJSON) != string(newFieldJSON), nil
}
