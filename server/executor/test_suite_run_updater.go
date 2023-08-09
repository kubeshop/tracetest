package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testsuite"
)

type TestSuiteRunUpdater interface {
	Update(context.Context, testsuite.TestSuiteRun) error
}

type CompositeTransactionUpdater struct {
	listeners []TestSuiteRunUpdater
}

func (u CompositeTransactionUpdater) Add(l TestSuiteRunUpdater) CompositeTransactionUpdater {
	u.listeners = append(u.listeners, l)
	return u
}

var _ TestSuiteRunUpdater = CompositeTransactionUpdater{}

func (u CompositeTransactionUpdater) Update(ctx context.Context, run testsuite.TestSuiteRun) error {
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
	UpdateRun(context.Context, testsuite.TestSuiteRun) error
}

func NewDBTranasctionUpdater(repo transactionUpdater) TestSuiteRunUpdater {
	return dbTransactionUpdater{repo}
}

func (u dbTransactionUpdater) Update(ctx context.Context, run testsuite.TestSuiteRun) error {
	return u.repo.UpdateRun(ctx, run)
}

type subscriptionTransactionUpdater struct {
	manager *subscription.Manager
}

func NewSubscriptionTransactionUpdater(manager *subscription.Manager) TestSuiteRunUpdater {
	return subscriptionTransactionUpdater{manager}
}

func (u subscriptionTransactionUpdater) Update(ctx context.Context, run testsuite.TestSuiteRun) error {
	u.manager.PublishUpdate(subscription.Message{
		ResourceID: run.ResourceID(),
		Type:       "result_update",
		Content:    run,
	})

	return nil
}
