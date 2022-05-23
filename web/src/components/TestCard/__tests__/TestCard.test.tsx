import {fireEvent, render, waitFor} from '@testing-library/react';
import TestMock from '../../../models/__mocks__/Test.mock';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import TestCard from '../TestCard';

test('TestCard', async () => {
  const onDelete = jest.fn();
  const onRunTest = jest.fn();
  const onClick = jest.fn();

<<<<<<< Updated upstream
  const test = TestMock.model();

  const {container, getByTestId} = render(
=======
  const {getByTestId, getByText} = render(
>>>>>>> Stashed changes
    <ReduxWrapperProvider>
      <TestCard onDelete={onDelete} onRunTest={onRunTest} test={test} onClick={onClick} />
    </ReduxWrapperProvider>
  );
<<<<<<< Updated upstream

  fireEvent(
    getByTestId(`test-actions-button-${test.id}`),
    new MouseEvent('click', {
      bubbles: true,
      cancelable: true,
    })
  );
  await waitFor(() => getByTestId('test-card-delete'));
  expect(container).toMatchSnapshot();
=======
  fireEvent(getByTestId(`test-actions-button-${TestingModels.test.testId}`), TestingModels.mouseEvent);
  await waitFor(() => getByTestId('delete'));
  expect(getByText('GET')).toBeTruthy();
  expect(getByText('Run Test')).toBeTruthy();
>>>>>>> Stashed changes
});
