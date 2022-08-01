import {render, waitFor} from 'test-utils';
import Layout from '../index';

test('Layout', async () => {
  const {getByText, getByTestId} = render(
    <Layout>
      <h2>This</h2>
    </Layout>
  );
  await waitFor(() => getByTestId('github-link'));

  expect(getByText('This')).toBeTruthy();
});
