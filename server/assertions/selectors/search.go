package selectors

import (
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

func filterSpans(rootSpan model.Span, spanSelector SpanSelector) []model.Span {
	filteredSpans := make([]model.Span, 0)
	traverseTree(rootSpan, func(span model.Span) {
		if spanSelector.MatchesFilters(span) {
			if spanSelector.ChildSelector != nil {
				childFilteredSpans := filterSpans(span, *spanSelector.ChildSelector)
				// filteredSpans include the parent span, so we have to remove it
				childFilteredSpans = removeSpanFromList(childFilteredSpans, span.ID)
				filteredSpans = append(filteredSpans, childFilteredSpans...)
			} else {
				filteredSpans = append(filteredSpans, span)
			}
		}
	})

	uniqueSpans := filterDuplicated(filteredSpans)

	if spanSelector.PseudoClass != nil {
		return spanSelector.PseudoClass.Filter(uniqueSpans)
	}

	return uniqueSpans
}

func traverseTree(rootNode model.Span, fn func(model.Span)) {
	// FIX: don't use recursion to prevent stackoverflow errors on huge traces
	fn(rootNode)
	for i := range rootNode.Children {
		child := rootNode.Children[i]
		traverseTree(*child, fn)
	}
}

func filterDuplicated(spans []model.Span) []model.Span {
	existingSpans := make(map[trace.SpanID]bool, 0)
	uniqueSpans := make([]model.Span, 0)
	for _, span := range spans {
		if _, exists := existingSpans[span.ID]; !exists {
			uniqueSpans = append(uniqueSpans, span)
			existingSpans[span.ID] = true
		}
	}

	return uniqueSpans
}

func removeSpanFromList(spans []model.Span, id trace.SpanID) []model.Span {
	idString := id.String()
	list := make([]model.Span, 0, len(spans))
	for _, span := range spans {
		if span.ID.String() != idString {
			list = append(list, span)
		}
	}

	return list
}
