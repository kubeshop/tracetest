package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
)

type RunUpdater interface {
	Update(context.Context, model.Run) error
}

type CompositeUpdater struct {
	listeners []RunUpdater
}

func (u CompositeUpdater) Add(l RunUpdater) CompositeUpdater {
	u.listeners = append(u.listeners, l)
	return u
}

var _ RunUpdater = CompositeUpdater{}

func (u CompositeUpdater) Update(ctx context.Context, run model.Run) error {
	for _, l := range u.listeners {
		if err := l.Update(ctx, run); err != nil {
			return fmt.Errorf("composite updating error: %w", err)
		}
	}

	return nil
}

type dbUpdater struct {
	repo model.RunRepository
}

func NewDBUpdater(repo model.RunRepository) RunUpdater {
	return dbUpdater{repo}
}

func (u dbUpdater) Update(ctx context.Context, run model.Run) error {
	return u.repo.UpdateRun(ctx, run)
}

type subscriptionUpdater struct {
	manager *subscription.Manager
}

func NewSubscriptionUpdater(manager *subscription.Manager) RunUpdater {
	return subscriptionUpdater{manager}
}

func (u subscriptionUpdater) Update(ctx context.Context, run model.Run) error {
	u.manager.PublishUpdate(subscription.Message{
		ResourceID: run.ResourceID(),
		Type:       "result_update",
		Content:    run,
	})

	return nil
}
