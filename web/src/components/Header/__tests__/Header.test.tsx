import {render, waitFor} from 'test-utils';
import Header from '../index';

test('Header', async () => {
  const {getByTestId} = render(<Header />);
  await waitFor(() => getByTestId('github-link'));
});
