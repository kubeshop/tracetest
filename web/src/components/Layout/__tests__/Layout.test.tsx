import {render} from 'test-utils';
import Layout from '../Layout';

it('Layout', async () => {
  const {getByText, getByTestId} = render(
    <Layout>
      <h2>This</h2>
    </Layout>
  );

  expect(getByTestId('menu-link')).toBeInTheDocument();
  expect(getByText('This')).toBeTruthy();
});
