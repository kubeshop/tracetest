import {fireEvent, render, waitFor} from '@testing-library/react';
import TestMock from '../../../models/__mocks__/Test.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import TestCard from '../TestCard';

test('TestCard', async () => {
  const onDelete = jest.fn();
  const onRunTest = jest.fn();
  const onClick = jest.fn();

  const test = TestMock.model();

  const {getByTestId, getByText} = render(
    <ReduxWrapperProvider>
      <TestCard onDelete={onDelete} onRunTest={onRunTest} test={test} onClick={onClick} />
    </ReduxWrapperProvider>
  );

  const mouseEvent = new MouseEvent('click', {bubbles: true});
  fireEvent(getByTestId(`test-actions-button-${test.id}`), mouseEvent);
  await waitFor(() => getByTestId('test-card-delete'));
  expect(getByText('Run Test')).toBeTruthy();
});
