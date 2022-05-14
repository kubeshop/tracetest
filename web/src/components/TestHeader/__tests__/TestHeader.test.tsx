import {render} from '@testing-library/react';
import TestHeader from '../TestHeader';
import {TestState} from '../../../constants/TestRunResult.constants';
import {HTTP_METHOD} from '../../../constants/Common.constants';

test('SpanAttributesTable', () => {
  const result = render(
    <TestHeader
      onBack={jest.fn()}
      test={{
        assertions: [],
        description: '',
        lastTestResult: undefined,
        name: '',
        serviceUnderTest: {
          id: '',
          request: {
            auth: undefined,
            body: '',
            certificate: undefined,
            headers: undefined,
            method: HTTP_METHOD.GET,
            proxy: undefined,
            url: '',
          },
        },
        testId: '',
      }}
      testState={TestState.CREATED}
    />
  );
  expect(result.container).toMatchSnapshot();
});
