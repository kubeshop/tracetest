import {render} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import Header from '../index';

test('Header', () => {
  const result = render(
    <MemoryRouter>
      <Header />
    </MemoryRouter>
  );
  expect(result.container).toMatchSnapshot();
});
