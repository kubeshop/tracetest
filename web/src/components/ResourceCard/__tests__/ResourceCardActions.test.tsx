import faker from '@faker-js/faker';
import {render, waitFor} from 'test-utils';
import ResourceCardActions from '../ResourceCardActions';

test('ResourceCardActions', async () => {
  const onEdit = jest.fn();
  const shouldEdit = true;
  const onDelete = jest.fn();
  const testId = faker.datatype.uuid();

  const {getByTestId} = render(
    <ResourceCardActions
      shouldEdit={shouldEdit}
      onEdit={onEdit}
      onDelete={onDelete}
      id={testId}
      onDuplicate={jest.fn()}
    />
  );
  await waitFor(() => getByTestId(`test-actions-button-${testId}`));
});
