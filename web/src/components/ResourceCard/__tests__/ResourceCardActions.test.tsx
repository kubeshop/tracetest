import faker from '@faker-js/faker';
import {render, waitFor} from 'test-utils';
import ResourceCardActions from '../ResourceCardActions';

test('ResourceCardActions', async () => {
  const onDelete = jest.fn();
  const testId = faker.datatype.uuid();

  const {getByTestId} = render(<ResourceCardActions onDelete={onDelete} id={testId} />);
  await waitFor(() => getByTestId(`test-actions-button-${testId}`));
});
