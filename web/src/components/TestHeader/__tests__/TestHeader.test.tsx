import {render} from 'test-utils';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestHeader from '../TestHeader';

test('SpanAttributesTable', () => {
  const test = TestMock.model();
  const {getByTestId} = render(
    <TestHeader description={test.description} id={test.id} onDelete={jest.fn} title={test.name} runButton={<div />} />
  );
  expect(getByTestId('test-details-name')).toBeInTheDocument();
});
