import {render} from '@testing-library/react';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {Provider} from 'react-redux';
import {CompareOperator} from '../../../constants/Operator.constants';
import {LOCATION_NAME} from '../../../constants/Span.constants';
import {SpanAttributeType} from '../../../constants/SpanAttribute.constants';
import {store} from '../../../redux/store';
import AssertionsTable from '../index';

test('AssetionsTable', () => {
  const result = render(
    <Provider store={store}>
      <AssertionsTable
        assertionResults={[
          {
            actualValue: '',
            comparisonValue: '',
            hasPassed: false,
            locationName: LOCATION_NAME.INSTRUMENTATION_LIBRARY,
            operator: CompareOperator.EQUALS,
            propertyName: '',
            spanAssertionId: '',
            spanId: '',
            valueType: SpanAttributeType.doubleValue,
          },
        ]}
        sort={1}
        assertion={{
          assertionId: '',
          selectors: [
            {
              locationName: LOCATION_NAME.INSTRUMENTATION_LIBRARY,
              propertyName: '',
              value: '',
              valueType: '',
            },
          ],
          spanAssertions: [],
        }}
        span={{
          attributeList: [],
          attributes: {},
          duration: 0,
          endTimeUnixNano: '',
          instrumentationLibrary: {name: '', version: ''},
          kind: '',
          name: '',
          parentSpanId: '',
          signature: [],
          spanId: '',
          startTimeUnixNano: '',
          status: {code: ''},
          traceId: '',
          type: SemanticGroupNames.Http,
        }}
        resultId="546"
        testId="234"
      />
    </Provider>
  );
  expect(result.container).toMatchSnapshot();
});
