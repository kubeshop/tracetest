import {render} from '@testing-library/react';
import {TestState} from '../../../constants/TestRun.constants';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestHeader from '../TestHeader';

test('SpanAttributesTable', () => {
  const result = render(<TestHeader onBack={jest.fn()} test={TestMock.model()} testState={TestState.CREATED} />);
  expect(result.container).toMatchSnapshot();
});
