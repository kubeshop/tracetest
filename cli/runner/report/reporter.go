package report

import (
	"fmt"

	"github.com/pterm/pterm"
)

type RunGroup struct {
	ID      string
	Summary RunGroupSummary
}

type RunGroupSummary struct {
	Passed     int
	Failed     int
	InProgress int
}

type Reporter struct {
	runGroupID  string
	runGroupUrl string

	multi                 *pterm.MultiPrinter
	groupSpinner          *pterm.SpinnerPrinter
	passedTestSpinner     *pterm.SpinnerPrinter
	failedTestSpinner     *pterm.SpinnerPrinter
	inProgressTestSpinner *pterm.SpinnerPrinter
}

func NewReporter(runGroupID, runGroupURL string) *Reporter {
	reporter := &Reporter{
		runGroupID:  runGroupID,
		runGroupUrl: runGroupURL,
		multi:       &pterm.DefaultMultiPrinter,
	}

	reporter.groupSpinner, _ = pterm.DefaultSpinner.WithWriter(reporter.multi.NewWriter()).Start(fmt.Sprintf(`Parallel tests: %s in progress`, runGroupID))
	reporter.passedTestSpinner, _ = pterm.DefaultSpinner.WithWriter(reporter.multi.NewWriter()).WithStyle(pterm.FgGreen.ToStyle()).Start("0 tests passed")
	reporter.failedTestSpinner, _ = pterm.DefaultSpinner.WithWriter(reporter.multi.NewWriter()).WithStyle(pterm.FgRed.ToStyle()).Start("0 tests failed")
	reporter.inProgressTestSpinner, _ = pterm.DefaultSpinner.WithWriter(reporter.multi.NewWriter()).Start("0 tests in progress")

	return reporter
}

func (r *Reporter) SetRunGroup(runGroup RunGroup) {
	r.render(runGroup)
}

func (r *Reporter) Start() {
	r.multi.Start()
}

func (r *Reporter) Stop() {
	r.multi.Stop()
}

func (r *Reporter) render(runGroup RunGroup) {
	r.passedTestSpinner.UpdateText(fmt.Sprintf("%d tests passed", runGroup.Summary.Passed))
	r.failedTestSpinner.UpdateText(fmt.Sprintf("%d tests failed", runGroup.Summary.Failed))
	r.inProgressTestSpinner.UpdateText(fmt.Sprintf("%d tests in progress", runGroup.Summary.InProgress))

	if runGroup.Summary.InProgress == 0 {
		if runGroup.Summary.Failed > 0 {
			r.groupSpinner.Fail(fmt.Sprintf(`Parallel tests: %s failed`, r.runGroupID))
			return
		}

		r.groupSpinner.Success(fmt.Sprintf(`Parallel tests: %s succeed`, r.runGroupID))
	}
}
