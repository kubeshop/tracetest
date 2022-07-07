import faker from '@faker-js/faker';
import {render, waitFor} from 'test-utils';
import TestCardActions from '../TestCardActions';

test('TestCardActions', async () => {
  const onDelete = jest.fn();
  const onEdit = jest.fn();
  const testId = faker.datatype.uuid();

  const {getByTestId} = render(<TestCardActions onDelete={onDelete} onEdit={onEdit} testId={testId} />);
  await waitFor(() => getByTestId(`test-actions-button-${testId}`));
});
