import {render, waitFor} from '@testing-library/react';
import {TestingModels} from '../../../utils/TestingModels';
import TestCardActions from '../TestCardActions';

test('TestCardActions', async () => {
  const onDelete = jest.fn();
  const {getByTestId, container} = render(<TestCardActions onDelete={onDelete} testId={TestingModels.testId} />);
  await waitFor(() => getByTestId(`test-actions-button-${TestingModels.testId}`));
  expect(container).toMatchSnapshot();
  // fireEvent(getByTestId(`test-actions-button-${TestingModels.test.testId}`), TestingModels.mouseEvent);
});
