package selectors

import (
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

func filterSpans(rootSpan traces.Span, spanSelector spanSelector) []traces.Span {
	filteredSpans := make([]traces.Span, 0)
	traverseTree(rootSpan, func(span traces.Span) {
		if spanSelector.MatchesFilters(span) {
			if spanSelector.ChildSelector != nil {
				childFilteredSpans := filterSpans(span, *spanSelector.ChildSelector)
				filteredSpans = append(filteredSpans, childFilteredSpans...)
			} else {
				filteredSpans = append(filteredSpans, span)
			}
		}
	})

	uniqueSpans := filterDuplicated(filteredSpans)

	if spanSelector.PsedoClass != nil {
		return spanSelector.PsedoClass.Filter(uniqueSpans)
	}

	return uniqueSpans
}

func traverseTree(rootNode traces.Span, fn func(traces.Span)) {
	// FIX: don't use recursion to prevent stackoverflow errors on huge traces
	fn(rootNode)
	for i := range rootNode.Children {
		child := rootNode.Children[i]
		traverseTree(*child, fn)
	}
}

func filterDuplicated(spans []traces.Span) []traces.Span {
	existingSpans := make(map[trace.SpanID]bool, 0)
	uniqueSpans := make([]traces.Span, 0)
	for _, span := range spans {
		if _, exists := existingSpans[span.ID]; !exists {
			uniqueSpans = append(uniqueSpans, span)
			existingSpans[span.ID] = true
		}
	}

	return uniqueSpans
}
