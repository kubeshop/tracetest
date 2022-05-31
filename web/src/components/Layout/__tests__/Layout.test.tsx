import {render, waitFor} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import Layout from '../index';

test('Layout', async () => {
  const {getByText, getByTestId} = render(
    <MemoryRouter>
      <Layout>
        <h2>This</h2>
      </Layout>
    </MemoryRouter>
  );
  await waitFor(() => getByTestId('github-link'));
  expect(getByText('This')).toBeTruthy();
});
