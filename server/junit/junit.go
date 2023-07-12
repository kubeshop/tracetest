package junit

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/test"
)

func FromRunResult(t test.Test, run test.Run) ([]byte, error) {

	assertions := []assertion{}

	var testTotals, testFails, testErrs int

	if run.Results == nil {
		return nil, fmt.Errorf("run has no results")
	}
	run.Results.Results.ForEach(func(selector test.SpanQuery, results []test.AssertionResult) error {
		checks := []check{}
		var total, fails, errs int
		for _, res := range results {
			for _, sar := range res.Results {
				total++
				c := check{
					Check: string(res.Assertion),
				}
				if sar.CompareErr != nil {
					if errors.Is(sar.CompareErr, comparator.ErrNoMatch) {
						fails++
						c.Failure = &failure{
							Assertion:   string(res.Assertion),
							ActualValue: sar.ObservedValue,
						}
					} else {
						errs++
						c.Error = &jError{
							Type:    "error", //arbitrary hardcoded value
							Message: sar.CompareErr.Error(),
						}
					}
				}
				checks = append(checks, c)
			}
		}

		testTotals += total
		testFails += fails
		testErrs += errs

		assertions = append(assertions, assertion{
			Selector: string(selector),
			Total:    total,
			Failures: fails,
			Errors:   errs,
			// Skipped: 0, // we don't have skips yet
			Checks: checks,
		})

		return nil
	})

	r := report{
		TestName:   t.Name,
		Time:       run.ExecutionTime(),
		Total:      testTotals,
		Failures:   testFails,
		Errors:     testErrs,
		Assertions: assertions,
	}

	return xml.MarshalIndent(r, "", "	")
}

type report struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestName   string   `xml:"name,attr"`
	Total      int      `xml:"tests,attr"`
	Failures   int      `xml:"failures,attr"`
	Errors     int      `xml:"errors,attr"`
	Skipped    int      `xml:"skipped,attr"`
	Time       int      `xml:"time,attr"`
	Assertions []assertion
}

type assertion struct {
	XMLName  xml.Name `xml:"testsuite"`
	Selector string   `xml:"name,attr"`
	Total    int      `xml:"tests,attr"`
	Failures int      `xml:"failures,attr"`
	Errors   int      `xml:"errors,attr"`
	Skipped  int      `xml:"skipped,attr"`
	Checks   []check
}

type check struct {
	XMLName xml.Name `xml:"testcase"`
	Check   string   `xml:"name,attr"`
	Error   *jError  `xml:",omitempty"`
	Failure *failure `xml:",omitempty"`
}

type jError struct {
	XMLName xml.Name `xml:"error"`
	Type    string   `xml:"type,attr"`
	Message string   `xml:"message,attr"`
}

type failure struct {
	XMLName     xml.Name `xml:"failure"`
	Assertion   string   `xml:"type,attr"`
	ActualValue string   `xml:"message,attr"`
}
