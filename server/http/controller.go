package http

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions"
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/junit"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testdb"
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
	test, err := c.mappers.In.Test(in)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	return c.doCreateTest(ctx, test)
}

var errTestExists = errors.New("test already exists")

func (c *controller) doCreateTest(ctx context.Context, test model.Test) (openapi.ImplResponse, error) {
	// if they try to create a test with preset ID, we need to make sure that ID doesn't exists already
	if test.HasID() {
		exists, err := c.testDB.IDExists(ctx, test.ID)

		if err != nil {
			return handleDBError(err), err
		}

		if exists {
			err := fmt.Errorf(`cannot create test with ID "%s: %w`, test.ID, errTestExists)
			r := map[string]string{
				"error": err.Error(),
			}
			return openapi.Response(http.StatusBadRequest, r), err
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

func (c *controller) GetTestResultSelectedSpans(ctx context.Context, testID, runID, selectorQuery string) (openapi.ImplResponse, error) {
	selector, err := selectors.New(selectorQuery)
	if err != nil {
		return handleDBError(err), err
	}

	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

func (c *controller) GetTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) DeleteTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteRun(ctx, run)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

type paginated[T any] struct {
	items []T
	count int
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

	return openapi.Response(200, paginated[openapi.TestRun]{
		items: c.mappers.Out.Runs(runs.Items),
		count: runs.TotalCount,
	}), nil
}

func (c *controller) GetTests(ctx context.Context, take, skip int32, query string, sortBy string, sortDirection string) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	tests, err := c.testDB.GetTests(ctx, take, skip, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, paginated[openapi.Test]{
		items: c.mappers.Out.Tests(tests.Items),
		count: tests.TotalCount,
	}), nil
}

func (c *controller) RerunTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

	metadata := metadata(runInformation.Metadata)
	environment := environment(ctx, c.testDB, runInformation.EnvironmentId)

	run := c.runner.Run(ctx, test, metadata, environment)

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

	newDefinition, err := c.mappers.In.Definition(def)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

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
	updated, err := c.mappers.In.Test(in)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	return c.doUpdateTest(ctx, id.ID(testID), updated)
}

func (c *controller) doUpdateTest(ctx context.Context, testID id.ID, updated model.Test) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, testID)
	if err != nil {
		return handleDBError(err), err
	}

	updated.Version = test.Version
	updated.ID = test.ID

	_, err = c.testDB.UpdateTest(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) DryRunAssertion(ctx context.Context, testID, runID string, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusUnprocessableEntity, fmt.Sprintf(`run "%s" has no trace associated`, runID)), nil
	}

	definition, err := c.mappers.In.Definition(def)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}
	results, allPassed := assertions.Assert(definition, *run.Trace)
	res := c.mappers.Out.Result(&model.RunResults{
		AllPassed: allPassed,
		Results:   results,
	})

	return openapi.Response(200, res), nil
}

func (c *controller) GetRunResultJUnit(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

	enc, err := yaml.Encode(yaml.File{
		Type: yaml.FileTypeTest,
		Spec: test,
	})
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	return openapi.Response(200, enc), nil
}

func (c controller) ExportTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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
	test, err := c.mappers.In.Test(exportedTest.Test)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	run, err := c.mappers.In.Run(exportedTest.Run)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

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

func (c *controller) ExecuteDefinition(ctx context.Context, testDefinition openapi.TextDefinition) (openapi.ImplResponse, error) {
	def, err := yaml.Decode([]byte(testDefinition.Content))
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	if test, err := def.Test(); err == nil {
		return c.executeTest(ctx, test.Model(), testDefinition.RunInformation)
	}

	return openapi.Response(http.StatusUnprocessableEntity, nil), nil
}

func metadata(in *map[string]string) model.RunMetadata {
	if in == nil {
		return nil
	}

	return model.RunMetadata(*in)
}

func environment(ctx context.Context, testDB model.Repository, environmentId string) model.Environment {
	if environmentId != "" {
		environment, err := testDB.GetEnvironment(ctx, environmentId)

		if err != nil {
			return model.Environment{}
		}

		return environment
	}

	return model.Environment{}
}

func (c *controller) executeTest(ctx context.Context, test model.Test, runInfo openapi.TestRunInformation) (openapi.ImplResponse, error) {
	// create or update test
	testID := test.ID
	resp, err := c.doCreateTest(ctx, test)
	if err != nil {
		if errors.Is(err, errTestExists) {
			resp, err := c.doUpdateTest(ctx, test.ID, test)
			if err != nil {
				return resp, err
			}

		} else {
			return resp, err
		}
	} else {
		testID = id.ID(resp.Body.(openapi.Test).Id)
	}

	// test ready, execute it
	resp, err = c.RunTest(ctx, testID.String(), runInfo)
	if err != nil {
		return resp, err
	}

	res := openapi.ExecuteDefinitionResponse{
		Id:    test.ID.String(),
		RunId: resp.Body.(openapi.TestRun).Id,
		Type:  string(yaml.FileTypeTest),
	}
	return openapi.Response(200, res), nil
}

// Environments

func (c *controller) CreateEnvironment(ctx context.Context, in openapi.Environment) (openapi.ImplResponse, error) {
	environment := c.mappers.In.Environment(in)

	if environment.ID != "" {
		exists, err := c.testDB.EnvironmentIDExists(ctx, environment.ID)

		if err != nil {
			return handleDBError(err), err
		}

		if exists {
			err := fmt.Errorf(`cannot create environment with ID "%s: %w`, environment.ID, errTestExists)
			r := map[string]string{
				"error": err.Error(),
			}
			return openapi.Response(http.StatusBadRequest, r), err
		}
	}

	environment, err := c.testDB.CreateEnvironment(ctx, environment)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, c.mappers.Out.Environment(environment)), nil
}

func (c *controller) DeleteEnvironment(ctx context.Context, environmentId string) (openapi.ImplResponse, error) {
	environment, err := c.testDB.GetEnvironment(ctx, environmentId)
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteEnvironment(ctx, environment)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) GetEnvironment(ctx context.Context, environmentId string) (openapi.ImplResponse, error) {
	environment, err := c.testDB.GetEnvironment(ctx, environmentId)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Environment(environment)), nil
}

func (c *controller) GetEnvironments(ctx context.Context, take, skip int32, query string, sortBy string, sortDirection string) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	environments, err := c.testDB.GetEnvironments(ctx, take, skip, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, paginated[openapi.Environment]{
		items: c.mappers.Out.Environments(environments.Items),
		count: environments.TotalCount,
	}), nil
}

func (c *controller) UpdateEnvironment(ctx context.Context, environmentId string, in openapi.Environment) (openapi.ImplResponse, error) {
	updated := c.mappers.In.Environment(in)

	environment, err := c.testDB.GetEnvironment(ctx, environmentId)
	if err != nil {
		return handleDBError(err), err
	}

	updated.ID = environment.ID

	_, err = c.testDB.UpdateEnvironment(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}
