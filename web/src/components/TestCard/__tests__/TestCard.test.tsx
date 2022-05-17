import {fireEvent, render, waitFor} from '@testing-library/react';
import {TestingModels} from '../../../utils/TestingModels';
import TestCard from '../TestCard';

const mouseEvent = new MouseEvent('click', {
  bubbles: true,
  cancelable: true,
});

test('TestCard', async () => {
  const onDelete = jest.fn();
  const onRunTest = jest.fn();
  const onClick = jest.fn();

  const {container, getByTestId} = render(
    <TestCard onDelete={onDelete} onRunTest={onRunTest} test={TestingModels.test} onClick={onClick} />
  );
  fireEvent(getByTestId('test-card'), mouseEvent);
  fireEvent(getByTestId('test-run-button'), mouseEvent);
  fireEvent(getByTestId(`test-actions-button-${TestingModels.test.testId}`), mouseEvent);
  await waitFor(() => getByTestId('delete'));
  expect(container).toMatchSnapshot();
});
