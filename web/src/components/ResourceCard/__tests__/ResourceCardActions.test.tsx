import faker from '@faker-js/faker';
import {render, waitFor} from 'test-utils';
import ResourceCardActions from '../ResourceCardActions';

test('ResourceCardActions', async () => {
  const onEdit = jest.fn();
  const canEdit = true;
  const onDelete = jest.fn();
  const testId = faker.datatype.uuid();

  const {getByTestId} = render(<ResourceCardActions canEdit={canEdit} onEdit={onEdit} onDelete={onDelete} id={testId} />);
  await waitFor(() => getByTestId(`test-actions-button-${testId}`));
});
