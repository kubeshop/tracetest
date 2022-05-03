package selectors

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/participle/v2"
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
	childSelector, err := createSpanSelectorFromParserSpanSelector(*parserSpanSelector.ChildSelector)
	if err != nil {
		return spanSelector{}, fmt.Errorf("could not create the child selector: %w", err)
	}

	pseudoClass, err := createPseudoClass(parserSpanSelector.PseudoClass)
	if err != nil {
		return spanSelector{}, err
	}

	filters := make([]filter, 0, len(parserSpanSelector.Filters))
	for _, parserFilter := range parserSpanSelector.Filters {
		filter := filter{
			Property: parserFilter.Property,
			// Operation: parserFilter.Operator,
			Value: Value{},
		}

		filters = append(filters, filter)
	}

	return spanSelector{
		Filters:       filters,
		PsedoClass:    pseudoClass,
		ChildSelector: &childSelector,
	}, nil
}

func createPseudoClass(parserPseudoClass parserPseudoClass) (pseudoClass, error) {
	switch parserPseudoClass.Type {
	case "nth_child":
		argument, err := strconv.Atoi(*parserPseudoClass.Value.String)
		if err != nil {
			return pseudoClass{}, fmt.Errorf("nth_child argument must be an integer")
		}

		return pseudoClass{
			Name:     parserPseudoClass.Type,
			Argument: Value{Type: ValueInt, Int: int64(argument)},
		}, nil
	}

	return pseudoClass{}, fmt.Errorf("unsupported pseudo class: %s", parserPseudoClass.Type)
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
