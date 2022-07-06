import {render} from 'test-utils';
import {TestState} from '../../../constants/TestRun.constants';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestHeader from '../TestHeader';

test('SpanAttributesTable', () => {
  const {getByTestId} = render(
    <TestHeader
      onBack={jest.fn()}
      showInfo={false}
      test={TestMock.model()}
      testState={TestState.CREATED}
      testVersion={1}
    />
  );
  expect(getByTestId('test-details-name')).toBeInTheDocument();
});
