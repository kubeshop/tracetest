import fetchMock from 'jest-fetch-mock';
import {store} from '../../store';
import {updateTestResult} from '../ResultList.slice';

describe('test ResultList slice', () => {
  it('updateTestResult', async () => {
    const resultId = '23049';
    fetchMock.mockResponse(JSON.stringify({}));
    const assertion = {assertionId: '', selectors: undefined, spanAssertions: undefined};
    await store.dispatch(
      updateTestResult({
        trace: {
          description: '',
          spans: [
            {
              attributeList: [],
              attributes: undefined,
              duration: 0,
              endTimeUnixNano: '',
              instrumentationLibrary: undefined,
              kind: '',
              name: '',
              parentSpanId: '',
              signature: [],
              spanId: '',
              startTimeUnixNano: '',
              status: {code: ''},
              traceId: '',
              type: undefined,
            },
          ],
        },
        resultId,
        test: {
          assertions: [assertion],
          description: '',
          lastTestResult: undefined,
          name: '',
          serviceUnderTest: {id: '', request: undefined},
          testId: '',
        },
      }) as any
    );

    expect((store.getState().resultList as any).resultListMap[resultId][0].assertion).toStrictEqual(assertion);
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
    expect((store.getState().resultList as any).resultListMap[resultId]).toStrictEqual([
      {assertion: {assertionId}, spanListAssertionResult: []},
    ]);
  });
});
