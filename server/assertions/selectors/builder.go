package selectors

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/traces"
)

func NewSelectorBuilder() (*SelectorBuilder, error) {
	parser, err := CreateParser()
	if err != nil {
		return nil, err
	}

	return &SelectorBuilder{
		parser: parser,
	}, nil
}

type SelectorBuilder struct {
	parser *participle.Parser
}

func (sb *SelectorBuilder) NewSelector(query string) (Selector, error) {
	parserSelector := ParserSelector{}
	err := sb.parser.ParseString("", query, &parserSelector)
	if err != nil {
		return Selector{}, fmt.Errorf("could not create selector from statement \"%s\": %w", query, err)
	}

	return createSelectorFromParserSelector(parserSelector)
}

func createSelectorFromParserSelector(parserSelector ParserSelector) (Selector, error) {
	selector := Selector{
		spanSelectors: make([]spanSelector, 0, len(parserSelector.SpanSelectors)),
	}
	for _, parserSpanSelector := range parserSelector.SpanSelectors {
		spanSelector, err := createSpanSelectorFromParserSpanSelector(parserSpanSelector)
		if err != nil {
			return Selector{}, err
		}

		selector.spanSelectors = append(selector.spanSelectors, spanSelector)
	}
	return Selector{}, nil
}

func createSpanSelectorFromParserSpanSelector(parserSpanSelector parserSpanSelector) (spanSelector, error) {
	var childSelector *spanSelector = nil
	if parserSpanSelector.ChildSelector != nil {
		newChildSelector, err := createSpanSelectorFromParserSpanSelector(*parserSpanSelector.ChildSelector)
		if err != nil {
			return spanSelector{}, fmt.Errorf("could not create the child selector: %w", err)
		}
		childSelector = &newChildSelector
	}

	pseudoClass, err := createPseudoClass(parserSpanSelector.PseudoClass)
	if err != nil {
		return spanSelector{}, err
	}

	filters := make([]filter, 0, len(parserSpanSelector.Filters))
	for _, parserFilter := range parserSpanSelector.Filters {
		operatorFunction, err := getOperatorFunction(parserFilter.Operator)
		if err != nil {
			return spanSelector{}, fmt.Errorf("could not create filter function: %w", err)
		}

		filter := filter{
			Property:  parserFilter.Property,
			Operation: operatorFunction,
			Value:     createValueFromParserValue(*parserFilter.Value),
		}

		filters = append(filters, filter)
	}

	return spanSelector{
		Filters:       filters,
		PsedoClass:    pseudoClass,
		ChildSelector: childSelector,
	}, nil
}

func getOperatorFunction(operator string) (filterFunction, error) {
	comparator, err := getComparatorFromOperator(operator)
	if err != nil {
		return nil, err
	}

	return func(span traces.Span, attribute string, value Value) error {
		attrValue := span.Attributes.Get(attribute)
		return comparator.Compare(value.AsString(), attrValue)
	}, nil
}

func getComparatorFromOperator(operator string) (comparator.Comparator, error) {
	registry := comparator.DefaultRegistry()
	comparator, err := registry.Get(operator)
	if err != nil {
		return nil, fmt.Errorf("Unsupported comparator %s: %w", operator, err)
	}

	return comparator, nil
}

func createPseudoClass(parserPseudoClass parserPseudoClass) (*pseudoClass, error) {
	switch parserPseudoClass.Type {
	case "nth_child":
		return &pseudoClass{
			Name:     parserPseudoClass.Type,
			Argument: Value{Type: ValueInt, Int: *parserPseudoClass.Value.Int},
		}, nil
	case "":
		// No pseudoClass
		return nil, nil
	}

	return nil, fmt.Errorf("unsupported pseudo class: %s", parserPseudoClass.Type)
}

func createValueFromParserValue(parserValue parserValue) Value {
	if parserValue.Boolean != nil {
		return Value{
			Type:    ValueBoolean,
			Boolean: *parserValue.Boolean,
		}
	}

	if parserValue.Float != nil {
		return Value{
			Type:  ValueFloat,
			Float: *parserValue.Float,
		}
	}

	if parserValue.Int != nil {
		return Value{
			Type: ValueInt,
			Int:  *parserValue.Int,
		}
	}

	if parserValue.String != nil {
		return Value{
			Type:   ValueString,
			String: *parserValue.String,
		}
	}

	return Value{
		Type: ValueNull,
	}
}
