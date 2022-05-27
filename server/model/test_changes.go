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

func BumpVersionIfDefinitionChanged(test Test, newDef Definition) (Test, error) {
	definitionHasChanged, err := testFieldHasChanged(test.Definition, newDef)
	if err != nil {
		return test, err
	}

	if definitionHasChanged {
		test.Version = test.Version + 1
	}

	return test, nil
}

func testHasChanged(oldTest Test, newTest Test) (bool, error) {
	definitionHasChanged, err := testFieldHasChanged(oldTest.Definition, newTest.Definition)
	if err != nil {
		return false, err
	}

	serviceUnderTestHasChanged, err := testFieldHasChanged(oldTest.ServiceUnderTest, newTest.ServiceUnderTest)
	if err != nil {
		return false, err
	}

	nameHasChanged := oldTest.Name != newTest.Name
	descriptionHasChanged := oldTest.Description != newTest.Description

	return definitionHasChanged || serviceUnderTestHasChanged || nameHasChanged || descriptionHasChanged, nil
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
