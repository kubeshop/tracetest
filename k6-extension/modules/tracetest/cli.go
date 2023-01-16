package tracetest

import (
	"fmt"

	"github.com/xoscar/xk6-tracetest-tracing/models"
	"github.com/xoscar/xk6-tracetest-tracing/utils"
)

func (t *Tracetest) getIsCliInstalled() bool {
	_, err := utils.RunCommand("tracetest", "version")

	return err == nil
}

func (t *Tracetest) exportTest(testID string) (string, error) {
	fileName := fmt.Sprintf("%s.yaml", testID)

	if utils.FileExists(fileName) {
		return fileName, nil
	}
	_, err := utils.RunCommand("tracetest", "test", "export", "--id", testID, "-o", fileName)
	return fileName, err
}

// TODO: use the traceId as input for the test
func (t *Tracetest) runTest(fileName, traceId string) (*models.TracetestRun, error) {
	res, err := utils.RunCommand("tracetest", "test", "run", "-d", fileName, "-o", "json")
	if err != nil {
		return nil, err
	}

	testRun := models.NewRun(res)
	err = utils.RemoveFile(fileName)

	// add link to the tracetest instance
	t.logger.Infoln(fmt.Sprintf("Test run path /test/%s/run/%s", testRun.Test.ID, testRun.TestRun.ID))
	return testRun, err
}

func (t *Tracetest) runFromDefinition(testDefinition, traceID string) (*models.TracetestRun, error) {
	fileName := fmt.Sprintf("%s.yaml", utils.RandHexStringRunes(8))
	err := utils.SaveFile(fileName, []byte(testDefinition))
	if err != nil {
		return nil, err
	}

	run, err := t.runTest(fileName, traceID)
	return run, err
}

func (t *Tracetest) runFromId(testID, traceID string) (*models.TracetestRun, error) {
	fileName, err := t.exportTest(testID)
	if err != nil {
		return nil, err
	}

	return t.runTest(fileName, traceID)
}
