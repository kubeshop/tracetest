import fetchMock from 'jest-fetch-mock';
import {HTTP_METHOD} from '../../../constants/Common.constants';
import {SemanticGroupNames} from '../../../constants/SemanticGroupNames.constants';
import {LOCATION_NAME} from '../../../constants/Span.constants';
import {IAssertion} from '../../../types/Assertion.types';
import {store} from '../../store';
import {updateTestResult} from '../ResultList.slice';

describe('test ResultList slice', () => {
  it('updateTestResult', async () => {
    const resultId = '23049';
    fetchMock.mockResponse(JSON.stringify({}));

    const assertion: IAssertion = {
      assertionId: '',
      selectors: [
        {
          propertyName: 'http.status_code',
          locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
          value: '200',
          valueType: 'intValue',
        },
      ],
      spanAssertions: [],
    };
    await store.dispatch(
      updateTestResult({
        trace: {
          description: '',
          spans: [
            {
              attributeList: [],
              attributes: {},
              duration: 0,
              endTimeUnixNano: '',
              kind: '',
              name: '',
              parentSpanId: '',
              signature: [],
              spanId: '',
              startTimeUnixNano: '',
              status: {code: ''},
              traceId: '',
              type: SemanticGroupNames.Http,
            },
          ],
        },
        resultId,
        test: {
          assertions: [assertion],
          description: '',
          lastTestResult: undefined,
          name: '',
          serviceUnderTest: {id: '', request: {url: 'http://localhost:3000', method: HTTP_METHOD.GET}},
          testId: '',
        },
      })
    );

    expect(store.getState().resultList.resultListMap[resultId][0].assertion).toStrictEqual(assertion);
  });

  it('dispatch resultList/replace', async () => {
    const resultId = '23049';
    const assertionId = '23e2389';
    store.dispatch({
      type: 'resultList/replace',
      payload: {
        assertionResult: [{assertionId, spanAssertionResults: []}],
        test: {assertions: [{assertionId}]},
        trace: {},
        resultId,
      },
    });

    expect(store.getState().resultList.resultListMap[resultId]).toStrictEqual([
      {assertion: {assertionId}, spanListAssertionResult: []},
    ]);
  });
});
