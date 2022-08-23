import faker from '@faker-js/faker';
import {render, waitFor} from 'test-utils';
import TestCardActions from '../TestCardActions';

test('TestCardActions', async () => {
  const onDelete = jest.fn();
  const testId = faker.datatype.uuid();

  const {getByTestId} = render(<TestCardActions onDelete={onDelete} testId={testId} />);
  await waitFor(() => getByTestId(`test-actions-button-${testId}`));
});
