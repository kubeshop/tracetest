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
});
