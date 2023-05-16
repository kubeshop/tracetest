package rules

import (
	"context"
	"fmt"
	"math"

	"github.com/kubeshop/tracetest/server/model"
)

type requiredAttributesRule struct {
	baseRule
	spanTypes          []string
	requiredAttributes []string
}

func NewRequiredAttributesRule(spanTypes []string, requiredAttributes []string) Rule {
	return &requiredAttributesRule{
		spanTypes:          spanTypes,
		requiredAttributes: requiredAttributes,
	}
}

func (r *requiredAttributesRule) Run(ctx context.Context, trace model.Trace) Result {
	spanResults := []SpanResult{}
	for _, targetType := range r.spanTypes {
		for _, span := range trace.Flat {
			if span.Attributes["tracetest.span.type"] != targetType {
				continue
			}

			missingAttributes := make([]string, 0)

			for _, requiredAttribute := range r.requiredAttributes {
				if _, attributeExists := span.Attributes[requiredAttribute]; !attributeExists {
					missingAttributes = append(missingAttributes, requiredAttribute)
				}
			}

			numberOfExpectedRequiredAttributes := float64(len(r.requiredAttributes))
			numberOfMissingRequiredAttributes := float64(len(missingAttributes))

			score := uint(math.Floor(((numberOfExpectedRequiredAttributes - numberOfMissingRequiredAttributes) / numberOfExpectedRequiredAttributes) * 100.0))

			errors := make([]string, 0, int(numberOfMissingRequiredAttributes))
			for _, missingAttribute := range missingAttributes {
				errors = append(errors, fmt.Sprintf(`Attribute "%s" is missing from span of type "%s"`, missingAttribute, targetType))
			}

			spanResults = append(spanResults, SpanResult{
				SpanID: span.ID.String(),
				Passed: numberOfMissingRequiredAttributes == 0,
				Score:  score,
				Errors: errors,
				Tips:   []string{},
			})
		}
	}

	var totalScore uint = 0
	var allPassed bool = true
	for _, spanResult := range spanResults {
		totalScore += spanResult.Score

		if !spanResult.Passed {
			allPassed = false
		}
	}

	var score uint

	if len(spanResults) == 0 {
		score = 100
	} else {
		score = totalScore / uint(len(spanResults))
	}

	return Result{
		Name:             r.name,
		Description:      r.description,
		Score:            score,
		Passed:           allPassed,
		NormalizedWeight: 0.0,
		SpansResults:     spanResults,
	}
}
