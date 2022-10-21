import {SemanticResourceAttributes, TraceTestAttributes} from 'constants/SpanAttribute.constants';
import {TSelector} from 'types/Common.types';
import {ISuggestion} from 'types/TestSpecs.types';
import {
  allSpansRule,
  createByAttributeAllSpansRule,
  createByAttributeDescendantsRule,
  selectedSpanRule,
  structuralPseudoClassRule,
} from './Rules';

const {TRACETEST_SPAN_TYPE} = TraceTestAttributes;
const {SERVICE_NAME} = SemanticResourceAttributes;

const rules = [
  allSpansRule,
  createByAttributeAllSpansRule(TRACETEST_SPAN_TYPE),
  createByAttributeAllSpansRule(SERVICE_NAME),
  structuralPseudoClassRule,
  createByAttributeDescendantsRule(TRACETEST_SPAN_TYPE),
  selectedSpanRule,
];

const SelectorSuggestionsService = {
  getSuggestions(
    selector: TSelector,
    matchedSpansId: string[],
    selectedSpanId: string,
    selectedSpanSelector: string,
    selectedParentSpanSelector: string
  ): ISuggestion[] {
    return rules
      .map(rule =>
        rule.match(selector, matchedSpansId, selectedSpanId, selectedSpanSelector, selectedParentSpanSelector)
      )
      .filter(suggestion => suggestion !== undefined) as ISuggestion[];
  },
};

export default SelectorSuggestionsService;
