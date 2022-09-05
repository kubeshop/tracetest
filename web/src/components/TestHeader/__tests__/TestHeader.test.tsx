import {render} from 'test-utils';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestHeader from '../TestHeader';

test('SpanAttributesTable', () => {
  const {getByTestId} = render(<TestHeader test={TestMock.model()} />);
  expect(getByTestId('test-details-name')).toBeInTheDocument();
});
