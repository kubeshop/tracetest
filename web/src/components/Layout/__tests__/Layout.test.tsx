import {render} from 'test-utils';
import Layout from '../Layout';

it('Layout', async () => {
  const {getByTestId} = render(<Layout />);

  expect(getByTestId('menu-link')).toBeInTheDocument();
});
