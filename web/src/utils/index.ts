import {COMPARE_OPERATOR} from 'types';

const flattenAttributesSelector = () =>
  `resourceSpans[].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0]},resource.attributes[].{key:key,value: value.*|[0]}]|[][]`;

export const filterBySpanId = (spanId: string = '') =>
  `resourceSpans[?instrumentationLibrarySpans[?spans[?starts_with(spanId,'${spanId}')]]] | [].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0],type:'span'},resource.attributes[].{key:key,value: value.*|[0],type:'resource'}]|[][]`;

export const getOperator = (op: COMPARE_OPERATOR) => {
  switch (op) {
    case COMPARE_OPERATOR.EQUALS:
      return 'eq';
    case COMPARE_OPERATOR.NOTEQUALS:
      return 'ne';
    case COMPARE_OPERATOR.GREATERTHAN:
      return 'gt';
    case COMPARE_OPERATOR.LESSTHAN:
      return 'lt';
      case COMPARE_OPERATOR.GREATOREQUALS:
      return 'gte';
    case COMPARE_OPERATOR.LESSOREQUAL:
      return 'lte';
    default:
      return 'eq';
  }
};
