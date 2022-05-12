import {act, render} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import Layout from '../index';

test('Layout', () => {
  act(() => {
    const result = render(
      <MemoryRouter>
        <Layout>
          <h2>This</h2>
        </Layout>
      </MemoryRouter>
    );
  });
  // const input = screen.getByText('This');
  // expect(result.container).toMatchSnapshot();
  // expect(input).toBeTruthy();
});
