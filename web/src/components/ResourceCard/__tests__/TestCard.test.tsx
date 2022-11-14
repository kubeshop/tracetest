import {fireEvent, render, waitFor} from 'test-utils';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestCard from '../TestCard';

test('TestCard', async () => {
  const onDelete = jest.fn();
  const onRunTest = jest.fn();
  const onClick = jest.fn();
  const test = TestMock.model();

  const {getByTestId, getByText} = render(
    <TestCard onDelete={onDelete} onRun={onRunTest} onViewAll={onClick} test={test} />
  );
  const mouseEvent = new MouseEvent('click', {bubbles: true});
  fireEvent(getByTestId(`test-actions-button-${test.id}`), mouseEvent);
  await waitFor(() => getByTestId('test-card-delete'));

  expect(getByText('Run')).toBeTruthy();
});
