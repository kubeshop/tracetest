import {search} from 'jmespath';
import {escapeString} from '../utils/Common';
import OperatorService from './Operator.service';
import {LOCATION_NAME} from '../constants/Span.constants';
import {getSpanValue} from './Span.service';
import {IResourceSpan} from '../types/Span.types';
import {ITrace} from '../types/Trace.types';
import {IAssertion, IAssertionResult, IItemSelector, ISpanAssertionResult} from '../types/Assertion.types';

const buildValueSelector = (comparisonValue: string, compareOperator: string, valueType: string) => {
  if (valueType === 'intValue') {
    return `to_number(value.${valueType}) ${compareOperator} \`${comparisonValue}\``;
  }

  if (compareOperator === 'contains') {
    return `contains(value.${valueType}, \`${comparisonValue}\`)`;
  }

  return `value.${valueType} ${compareOperator} \`${comparisonValue}\``;
};

const buildSelector = (locationName: LOCATION_NAME, conditions: string[], spanId?: string) => {
  switch (locationName) {
    case LOCATION_NAME.INSTRUMENTATION_LIBRARY:
    case LOCATION_NAME.RESOURCE_ATTRIBUTES:
      return `${conditions.map(cond => `resource.attributes[?${cond}]`).join(' && ')}`;
    case LOCATION_NAME.SPAN_ATTRIBUTES:
    case LOCATION_NAME.SPAN:
      return `instrumentationLibrarySpans[?spans[?${conditions.map(cond => `attributes[?${cond}]`).join(' && ')}]]`;
    case LOCATION_NAME.SPAN_ID:
      return `instrumentationLibrarySpans[?spans[?spanId == \`${spanId}\` ${conditions.length ? '&&' : ''} ${conditions
        .map(cond => `attributes[?${cond}]`)
        .join(' && ')}]]`;
    default:
      return '';
  }
};

const buildConditionArray = (itemSelectors: IItemSelector[]) => {
  const selectorsMap = itemSelectors.reduce<Record<string, string[]>>(
    (acc, {locationName, propertyName, valueType, value}) => {
      const keySelector = ` key == \`${propertyName}\``;
      const valueSelector = buildValueSelector(value, '==', valueType);
      const condition = ` ${[keySelector, valueSelector]!.join(' && ')}`;

      return {
        ...acc,
        [locationName]: acc[locationName] ? acc[locationName].concat([condition]) : [condition],
      };
    },
    {}
  );

  return Object.entries(selectorsMap).map(([locationName, conditionList]) =>
    buildSelector(locationName as LOCATION_NAME, conditionList)
  );
};

export const runSpanAssertionByResourceSpan = (
  span: IResourceSpan,
  spanId: string,
  {spanAssertions = []}: IAssertion
): Array<ISpanAssertionResult> => {
  const assertionTestResultArray = spanAssertions.map(spanAssertion => {
    const {comparisonValue, operator, valueType, locationName, propertyName} = spanAssertion;
    const valueSelector = buildValueSelector(comparisonValue, OperatorService.getOperatorSymbol(operator), valueType);

    let selector = `${buildSelector(locationName, [`key == \`${propertyName}\` && ${valueSelector}`])}`;

    if ([LOCATION_NAME.SPAN, LOCATION_NAME.SPAN_ATTRIBUTES].includes(locationName)) {
      selector = `${buildSelector(LOCATION_NAME.SPAN_ID, [`key == \`${propertyName}\` && ${valueSelector}`], spanId)}`;
    }

    const [passedSpan] = search(span, escapeString(selector));

    return {
      ...spanAssertion,
      spanId,
      hasPassed: Boolean(passedSpan),
      actualValue: getSpanValue(span, locationName, propertyName),
    };
  });

  return assertionTestResultArray;
};

export const runAssertionBySpanId = (
  spanId: string,
  trace: ITrace,
  assertion: IAssertion
): ISpanAssertionResult[] | undefined => {
  if (!assertion.selectors) return undefined;
  const conditionList = buildConditionArray(assertion.selectors);
  const itemSelector = `[? ${buildSelector(LOCATION_NAME.SPAN_ID, [], spanId)} && ${conditionList.join(' && ')}]`;

  const [span]: Array<IResourceSpan> = search(trace, escapeString(`resourceSpans|[]| ${itemSelector}`));

  if (!span) return undefined;

  return runSpanAssertionByResourceSpan(span, spanId, assertion);
};

export const runAssertionByTrace = (trace: ITrace, assertion: IAssertion): IAssertionResult => {
  if (!assertion?.selectors) return {assertion, spanListAssertionResult: []};

  const itemSelector = `[? ${buildConditionArray(assertion.selectors).join(' && ')}]`;
  const spanList: Array<IResourceSpan> = search(trace, escapeString(`resourceSpans|[]| ${itemSelector}`)) || [];

  return {
    assertion,
    spanListAssertionResult: spanList.map(span => ({
      span,
      resultList: span.instrumentationLibrarySpans.reduce<ISpanAssertionResult[]>(
        (resultList, instrumentationLibrary) =>
          resultList.concat(
            instrumentationLibrary.spans
              .map(({spanId}) => runSpanAssertionByResourceSpan(span, spanId, assertion))
              .flat()
          ),
        []
      ),
    })),
  };
};

export const getEffectedSpansCount = (trace: ITrace, selectors: IItemSelector[]) => {
  if (selectors.length === 0) return 0;

  const itemSelector = `[? ${buildConditionArray(selectors).join(' && ')}]`;
  const spanList: Array<IResourceSpan> = search(trace, escapeString(`resourceSpans|[]| ${itemSelector}`));

  return spanList.length;
};
