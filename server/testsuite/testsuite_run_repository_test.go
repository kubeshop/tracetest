package testsuite_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestWithName(t *testing.T, db test.Repository, name string) test.Test {
	t.Helper()
	test := test.Test{
		Name:        name,
		Description: "description",
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
			HTTP: &trigger.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
	}

	updated, err := db.Create(context.TODO(), test)
	if err != nil {
		panic(err)
	}
	return updated
}

func getRepos() (*testsuite.Repository, *testsuite.RunRepository, test.Repository) {
	db := testmock.CreateMigratedDatabase()

	testRepo := test.NewRepository(db)
	testRunRepo := test.NewRunRepository(db, test.NewCache("test"))

	transactionRepo := testsuite.NewRepository(db, testRepo)
	runRepo := testsuite.NewRunRepository(db, testRunRepo)

	return transactionRepo, runRepo, testRepo
}

func getTransaction(t *testing.T, transactionRepo *testsuite.Repository, testsRepo test.Repository) (testsuite.TestSuite, testSuiteFixture) {
	f := setupTestSuiteFixture(t, transactionRepo.DB())

	suite := testsuite.TestSuite{
		ID:          id.NewRandGenerator().ID(),
		Name:        "first test",
		Description: "description",
		StepIDs: []id.ID{
			f.t1.ID,
			f.t2.ID,
		},
	}

	_, err := transactionRepo.Create(context.TODO(), suite)
	require.NoError(t, err)

	suite, err = transactionRepo.GetAugmented(context.TODO(), suite.ID)
	require.NoError(t, err)

	return suite, f
}

func TestCreateTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	transactionObject, _ := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), transactionObject.NewRun())
	require.NoError(t, err)

	assert.Equal(t, tr.TestSuiteID, transactionObject.ID)
	assert.Equal(t, tr.State, testsuite.TestSuiteStateCreated)
	assert.Len(t, tr.Steps, 0)
}

func TestUpdateTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	transactionObject, fixture := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), transactionObject.NewRun())
	require.NoError(t, err)

	tr.State = testsuite.TestSuiteStateExecuting
	tr.Steps = []test.Run{fixture.testRun}
	err = transactionRunRepo.UpdateRun(context.TODO(), tr)
	require.NoError(t, err)

	updatedRun, err := transactionRunRepo.GetTestSuiteRun(context.TODO(), transactionObject.ID, tr.ID)
	require.NoError(t, err)

	assert.Equal(t, tr.TestSuiteID, transactionObject.ID)
	assert.Equal(t, testsuite.TestSuiteStateExecuting, updatedRun.State)
	assert.Len(t, tr.Steps, 1)
}

func TestDeleteTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()
	suite, _ := getTransaction(t, transactionRepo, testsRepo)

	tr, err := transactionRunRepo.CreateRun(context.TODO(), suite.NewRun())
	require.NoError(t, err)

	err = transactionRunRepo.DeleteTestSuiteRun(context.TODO(), tr)
	require.NoError(t, err)

	_, err = transactionRunRepo.GetTestSuiteRun(context.TODO(), suite.ID, tr.ID)
	require.ErrorContains(t, err, "no rows in result set")
}

func createTransaction(t *testing.T, repo *testsuite.Repository, suite testsuite.TestSuite) testsuite.TestSuite {
	one := 1
	suite.ID = id.GenerateID()
	suite.Version = &one
	for _, step := range suite.Steps {
		suite.StepIDs = append(suite.StepIDs, step.ID)
	}
	_, err := repo.Create(context.TODO(), suite)
	require.NoError(t, err)

	suite, err = repo.GetAugmented(context.TODO(), suite.ID)
	require.NoError(t, err)

	return suite
}

func TestListTransactionRun(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()

	t1 := createTransaction(t, transactionRepo, testsuite.TestSuite{
		Name:        "first test",
		Description: "description",
		Steps: []test.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	t2 := createTransaction(t, transactionRepo, testsuite.TestSuite{
		Name:        "second testsuite",
		Description: "description",
		Steps: []test.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	run1, err := transactionRunRepo.CreateRun(context.TODO(), t1.NewRun())
	require.NoError(t, err)

	run2, err := transactionRunRepo.CreateRun(context.TODO(), t1.NewRun())
	require.NoError(t, err)

	_, err = transactionRunRepo.CreateRun(context.TODO(), t2.NewRun())
	require.NoError(t, err)

	runs, err := transactionRunRepo.GetTestSuiteRuns(context.TODO(), t1.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 2)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)
}

func TestBug(t *testing.T) {
	transactionRepo, transactionRunRepo, testsRepo := getRepos()

	ctx := context.TODO()

	suite := createTransaction(t, transactionRepo, testsuite.TestSuite{
		Name:        "first test",
		Description: "description",
		Steps: []test.Test{
			createTestWithName(t, testsRepo, "first step"),
			createTestWithName(t, testsRepo, "second step"),
		},
	})

	run1, err := transactionRunRepo.CreateRun(ctx, suite.NewRun())
	require.NoError(t, err)

	run2, err := transactionRunRepo.CreateRun(ctx, suite.NewRun())
	require.NoError(t, err)

	runs, err := transactionRunRepo.GetTestSuiteRuns(ctx, suite.ID, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, runs[0].ID, run2.ID)
	assert.Equal(t, runs[1].ID, run1.ID)

	suite.Name = "another thing"
	_, err = transactionRepo.Update(ctx, suite)
	require.NoError(t, err)

	newTransaction, err := transactionRepo.GetAugmented(context.TODO(), suite.ID)
	require.NoError(t, err)

	run3, err := transactionRunRepo.CreateRun(ctx, suite.NewRun())
	require.NoError(t, err)

	runs, err = transactionRunRepo.GetTestSuiteRuns(ctx, newTransaction.ID, 20, 0)
	require.NoError(t, err)

	assert.Len(t, runs, 3)
	assert.Equal(t, runs[0].ID, run3.ID)
	assert.Equal(t, runs[1].ID, run2.ID)
	assert.Equal(t, runs[2].ID, run1.ID)

}
