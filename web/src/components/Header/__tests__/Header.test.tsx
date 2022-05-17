import {render, waitFor} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import Header from '../index';

test('Header', async () => {
  const {container, getByTestId} = render(
    <MemoryRouter>
      <Header />
    </MemoryRouter>
  );

  await waitFor(() => getByTestId('github-link'));

  expect(container).toMatchSnapshot();
});
