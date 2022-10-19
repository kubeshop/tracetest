import {PseudoSelector} from 'constants/Operator.constants';
import {SemanticResourceAttributes, TraceTestAttributes} from 'constants/SpanAttribute.constants';
import {TSelector} from 'types/Common.types';
import {ISuggestion} from 'types/TestSpecs.types';
import {allSpansRule, createByAttributeRule} from './Rules';

const {TRACETEST_SPAN_TYPE} = TraceTestAttributes;
const {SERVICE_NAME} = SemanticResourceAttributes;

const rules = [
  allSpansRule,
  createByAttributeRule(TRACETEST_SPAN_TYPE),
  createByAttributeRule(SERVICE_NAME),
  createByAttributeRule(TRACETEST_SPAN_TYPE, PseudoSelector.FIRST),
];

const SelectorSuggestionsService = {
  getSuggestions(selector: TSelector): ISuggestion[] {
    return rules.map(rule => rule.match(selector)).filter(suggestion => suggestion !== undefined) as ISuggestion[];
  },
};

export default SelectorSuggestionsService;
