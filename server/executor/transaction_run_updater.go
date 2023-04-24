package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tests"
)

type TransactionRunUpdater interface {
	Update(context.Context, tests.TransactionRun) error
}

type CompositeTransactionUpdater struct {
	listeners []TransactionRunUpdater
}

func (u CompositeTransactionUpdater) Add(l TransactionRunUpdater) CompositeTransactionUpdater {
	u.listeners = append(u.listeners, l)
	return u
}

var _ TransactionRunUpdater = CompositeTransactionUpdater{}

func (u CompositeTransactionUpdater) Update(ctx context.Context, run tests.TransactionRun) error {
	for _, l := range u.listeners {
		if err := l.Update(ctx, run); err != nil {
			return fmt.Errorf("composite updating error: %w", err)
		}
	}

	return nil
}

type dbTransactionUpdater struct {
	repo transactionUpdater
}

type transactionUpdater interface {
	UpdateRun(context.Context, tests.TransactionRun) error
}

func NewDBTranasctionUpdater(repo transactionUpdater) TransactionRunUpdater {
	return dbTransactionUpdater{repo}
}

func (u dbTransactionUpdater) Update(ctx context.Context, run tests.TransactionRun) error {
	return u.repo.UpdateRun(ctx, run)
}

type subscriptionTransactionUpdater struct {
	manager *subscription.Manager
}

func NewSubscriptionTransactionUpdater(manager *subscription.Manager) TransactionRunUpdater {
	return subscriptionTransactionUpdater{manager}
}

func (u subscriptionTransactionUpdater) Update(ctx context.Context, run tests.TransactionRun) error {
	u.manager.PublishUpdate(subscription.Message{
		ResourceID: run.ResourceID(),
		Type:       "result_update",
		Content:    run,
	})

	return nil
}
