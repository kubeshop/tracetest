package http

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/kubeshop/tracetest/server/assertions"
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/encoding/yaml/conversion"
	"github.com/kubeshop/tracetest/server/encoding/yaml/definition"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/junit"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testdb"
	"gopkg.in/yaml.v3"
)

var IDGen = id.NewRandGenerator()

type controller struct {
	testDB          model.Repository
	runner          executor.Runner
	assertionRunner executor.AssertionRunner
	mappers         mappings.Mappings
}

func NewController(
	testDB model.Repository,
	runner executor.Runner,
	assertionRunner executor.AssertionRunner,
	mappers mappings.Mappings,
) openapi.ApiApiServicer {
	return &controller{
		testDB:          testDB,
		runner:          runner,
		assertionRunner: assertionRunner,
		mappers:         mappers,
	}
}

func handleDBError(err error) openapi.ImplResponse {
	switch {
	case errors.Is(testdb.ErrNotFound, err):
		return openapi.Response(http.StatusNotFound, err.Error())
	default:
		return openapi.Response(http.StatusInternalServerError, err.Error())
	}
}

func (c *controller) CreateTest(ctx context.Context, in openapi.Test) (openapi.ImplResponse, error) {
	test := c.mappers.In.Test(in)

	// if they try to create a test with preset ID, we need to make sure that ID doesn't exists already
	if test.HasID() {
		exists, err := c.testDB.IDExists(ctx, test.ID)

		if err != nil {
			return handleDBError(err), err
		}

		if exists {
			r := map[string]string{
				"error": fmt.Sprintf(`test with ID "%s" already exists. try updating instead`, test.ID),
			}
			return openapi.Response(http.StatusBadRequest, r), nil
		}
	}

	test, err := c.testDB.CreateTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c *controller) DeleteTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) GetTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c *controller) GetTestSpecs(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Specs(test.Specs)), nil
}

func (c *controller) GetTestResultSelectedSpans(ctx context.Context, _ string, runID int32, selectorQuery string) (openapi.ImplResponse, error) {
	selector, err := selectors.New(selectorQuery)
	if err != nil {
		return handleDBError(err), err
	}

	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusInternalServerError, "trace not available"), nil
	}

	selectedSpans := selector.Filter(*run.Trace)
	selectedSpanIds := make([]string, len(selectedSpans))

	for i, span := range selectedSpans {
		selectedSpanIds[i] = hex.EncodeToString(span.ID[:])
	}

	return openapi.Response(http.StatusOK, selectedSpanIds), nil
}

func (c *controller) GetTestRun(ctx context.Context, _ string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) DeleteTestRun(ctx context.Context, _ string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteRun(ctx, run)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) GetTestRuns(ctx context.Context, testID string, take, skip int32) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	runs, err := c.testDB.GetTestRuns(ctx, test, take, skip)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Runs(runs)), nil
}

func (c *controller) GetTests(ctx context.Context, take, skip int32, query string) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	tests, err := c.testDB.GetTests(ctx, take, skip, query)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Tests(tests)), nil
}

func (c *controller) RerunTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {

	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	newTestRun, err := c.testDB.CreateRun(ctx, test, run.Copy())
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	newTestRun = newTestRun.SuccessfullyPolledTraces(run.Trace)
	err = c.testDB.UpdateRun(ctx, newTestRun)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	assertionRequest := executor.AssertionRequest{
		Test: test,
		Run:  newTestRun,
	}

	c.assertionRunner.RunAssertions(ctx, assertionRequest)

	return openapi.Response(http.StatusOK, c.mappers.Out.Run(&newTestRun)), nil
}

func (c *controller) RunTest(ctx context.Context, testID string, runInformation openapi.TestRunInformation) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	metadata := model.RunMetadata{}
	if runInformation.Metadata != nil {
		metadata = model.RunMetadata(*runInformation.Metadata)
	}
	run := c.runner.Run(ctx, test, metadata)

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) SetTestSpecs(ctx context.Context, testID string, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	if err := c.mappers.In.ValidateDefinition(def); err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	newDefinition := c.mappers.In.Definition(def)

	newTest, err := model.BumpVersionIfDefinitionChanged(test, newDefinition)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	newTest.Specs = newDefinition

	newTest, err = c.testDB.UpdateTest(ctx, newTest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) UpdateTest(ctx context.Context, testID string, in openapi.Test) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	updated := c.mappers.In.Test(in)
	updated.Version = test.Version
	updated.ID = test.ID
	updated.ReferenceRun = nil

	_, err = c.testDB.UpdateTest(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) DryRunAssertion(ctx context.Context, _ string, runID int32, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusUnprocessableEntity, fmt.Sprintf(`run "%d" has no trace associated`, runID)), nil
	}

	results, allPassed := assertions.Assert(c.mappers.In.Definition(def), *run.Trace)
	res := c.mappers.Out.Result(&model.RunResults{
		AllPassed: allPassed,
		Results:   results,
	})

	return openapi.Response(200, res), nil
}

func (c *controller) GetRunResultJUnit(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), run.TestVersion)
	if err != nil {
		return handleDBError(err), err
	}

	res, err := junit.FromRunResult(test, run)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, res), nil
}

func (c controller) GetTestVersion(ctx context.Context, testID string, version int32) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), int(version))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c controller) GetTestVersionDefinitionFile(ctx context.Context, testID string, version int32) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), int(version))
	if err != nil {
		return handleDBError(err), err
	}

	res, err := conversion.GetYamlFileFromDefinition(c.mappers.Out.TestDefinitionFile(test))
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	return openapi.Response(200, res), nil
}

func (c controller) ExportTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), run.TestVersion)
	if err != nil {
		return handleDBError(err), err
	}

	response := openapi.ExportedTestInformation{
		Test: c.mappers.Out.Test(test),
		Run:  c.mappers.Out.Run(&run),
	}

	return openapi.Response(http.StatusOK, response), nil
}

func (c controller) ImportTestRun(ctx context.Context, exportedTest openapi.ExportedTestInformation) (openapi.ImplResponse, error) {
	test := c.mappers.In.Test(exportedTest.Test)
	run := c.mappers.In.Run(exportedTest.Run)

	createdTest, err := c.testDB.CreateTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	createdRun, err := c.testDB.CreateRun(ctx, createdTest, *run)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	createdRun.State = run.State

	err = c.testDB.UpdateRun(ctx, createdRun)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	response := openapi.ExportedTestInformation{
		Test: c.mappers.Out.Test(createdTest),
		Run:  c.mappers.Out.Run(&createdRun),
	}

	return openapi.Response(http.StatusOK, response), nil
}

func (c *controller) CreateTestFromDefinition(ctx context.Context, testDefinition openapi.TextDefinition) (openapi.ImplResponse, error) {
	var definitionObject definition.Test
	err := yaml.Unmarshal([]byte(testDefinition.Content), &definitionObject)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	openapiObject, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definitionObject)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	return c.CreateTest(ctx, openapiObject)
}

func (c *controller) UpdateTestFromDefinition(ctx context.Context, testId string, testDefinition openapi.TextDefinition) (openapi.ImplResponse, error) {
	var definitionObject definition.Test
	err := yaml.Unmarshal([]byte(testDefinition.Content), &definitionObject)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	openapiObject, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definitionObject)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	response, err := c.UpdateTest(ctx, testId, openapiObject)
	if err != nil {
		return response, err
	}

	return openapi.Response(http.StatusOK, openapiObject), nil
}

func (c *controller) RunShortUrl(ctx context.Context, shortID string) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRunByShortID(ctx, shortID)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.ImplResponse{
		Body: fmt.Sprintf("/test/%s/run/%d", run.TestID, run.ID),
	}, nil
}
