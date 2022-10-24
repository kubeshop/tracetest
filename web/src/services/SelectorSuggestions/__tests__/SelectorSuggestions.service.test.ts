import {TSelector} from 'types/Common.types';
import SelectorSuggestionsService from '../SelectorSuggestions.service';

const matchedSpansId: string[] = [];
const selectedSpanId = '';
const selectedSpanSelector = '';
const selectedParentSpanSelector = '';

describe('SelectorSuggestionsService', () => {
  it('should get All spans suggestion', () => {
    const selector: TSelector = {
      query: 'span[attribute="value"]',
      structure: [
        {
          filters: [{operator: '=', property: 'attribute', value: 'value'}],
        },
      ],
    };
    const suggestions = SelectorSuggestionsService.getSuggestions(
      selector,
      matchedSpansId,
      selectedSpanId,
      selectedSpanSelector,
      selectedParentSpanSelector
    );

    expect(suggestions).toContainEqual({query: '', title: 'All spans'});
  });

  it('should get TRACETEST_SPAN_TYPE suggestion', () => {
    const selector: TSelector = {
      query: 'span[tracetest.span.type="general" attribute="value"]',
      structure: [
        {
          filters: [
            {operator: '=', property: 'tracetest.span.type', value: 'general'},
            {operator: '=', property: 'attribute', value: 'value'},
          ],
        },
      ],
    };
    const suggestions = SelectorSuggestionsService.getSuggestions(
      selector,
      matchedSpansId,
      selectedSpanId,
      selectedSpanSelector,
      selectedParentSpanSelector
    );

    expect(suggestions).toContainEqual({query: 'span[tracetest.span.type="general"]', title: 'All general spans'});
  });

  it('should get SERVICE_NAME suggestion', () => {
    const selector: TSelector = {
      query: 'span[attribute="value" service.name="cart-api"]',
      structure: [
        {
          filters: [
            {operator: '=', property: 'attribute', value: 'value'},
            {operator: '=', property: 'service.name', value: 'cart-api'},
          ],
        },
      ],
    };
    const suggestions = SelectorSuggestionsService.getSuggestions(
      selector,
      matchedSpansId,
      selectedSpanId,
      selectedSpanSelector,
      selectedParentSpanSelector
    );

    expect(suggestions).toContainEqual({query: 'span[service.name="cart-api"]', title: 'All cart-api spans'});
  });

  it('should get structural pseudo class suggestion', () => {
    const selector: TSelector = {
      query: 'span[tracetest.span.type="general" attribute="value"]',
      structure: [
        {
          filters: [
            {operator: '=', property: 'tracetest.span.type', value: 'general'},
            {operator: '=', property: 'attribute', value: 'value'},
          ],
        },
      ],
    };
    const suggestions = SelectorSuggestionsService.getSuggestions(
      selector,
      ['123', '456'],
      '123',
      selectedSpanSelector,
      selectedParentSpanSelector
    );

    expect(suggestions).toContainEqual({query: `${selector.query}:first`, title: '1st span in group'});
  });

  it('should get descendant suggestion', () => {
    const selector: TSelector = {
      query: 'span[tracetest.span.type="general" attribute="value"]',
      structure: [
        {
          filters: [
            {operator: '=', property: 'tracetest.span.type', value: 'general'},
            {operator: '=', property: 'attribute', value: 'value'},
          ],
        },
      ],
    };
    const suggestions = SelectorSuggestionsService.getSuggestions(
      selector,
      matchedSpansId,
      selectedSpanId,
      selectedSpanSelector,
      'span[parent="parent"]'
    );

    expect(suggestions).toContainEqual({
      query: `span[parent="parent"] span[tracetest.span.type="general"]`,
      title: 'general descendants',
    });
  });

  it('should get selected span suggestion', () => {
    const selector: TSelector = {
      query: 'span[tracetest.span.type="general" attribute="value"]',
      structure: [
        {
          filters: [
            {operator: '=', property: 'tracetest.span.type', value: 'general'},
            {operator: '=', property: 'attribute', value: 'value'},
          ],
        },
      ],
    };
    const suggestions = SelectorSuggestionsService.getSuggestions(
      selector,
      matchedSpansId,
      selectedSpanId,
      'span[name="selected-span"]',
      selectedParentSpanSelector
    );

    expect(suggestions).toContainEqual({
      query: 'span[name="selected-span"]',
      title: 'Current span',
    });
  });
});
