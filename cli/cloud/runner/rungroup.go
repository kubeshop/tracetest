package runner

import (
	"context"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type runGroup struct {
	openapiClient *openapi.APIClient
	logger        *zap.Logger
}

func RunGroup(openapiClient *openapi.APIClient) *runGroup {
	return &runGroup{
		logger:        cmdutil.GetLogger(),
		openapiClient: openapiClient,
	}
}

func (rg *runGroup) WaitForCompletion(ctx context.Context, runGroupID string) (openapi.RunGroup, error) {
	var (
		updatedResult openapi.RunGroup
		lastError     error
		wg            sync.WaitGroup
	)

	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: change to websockets
	go func() {
		for range ticker.C {
			req := rg.openapiClient.ApiApi.GetRunGroup(ctx, runGroupID)
			runGroup, _, err := req.Execute()

			// updatedResult = runGroup
			rg.logger.Debug("updated run group", zap.String("result", spew.Sdump(runGroup)))
			if err != nil {
				rg.logger.Debug("UpdateResult failed", zap.Error(err))
				lastError = err
				wg.Done()
				return
			}

			if runGroup.GetStatus() == "succeed" || runGroup.GetStatus() == "failed" {
				rg.logger.Debug("result is finished")
				updatedResult = *runGroup
				wg.Done()
				return
			}
			rg.logger.Debug("still waiting")
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.RunGroup{}, lastError
	}

	return updatedResult, nil
}
