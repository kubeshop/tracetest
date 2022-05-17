import {render, waitFor} from '@testing-library/react';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import Router from '../Router';

test('Router', async () => {
  const {container, getByTestId} = render(
    <ReduxWrapperProvider>
      <Router />
    </ReduxWrapperProvider>
  );

  await waitFor(() => getByTestId('github-link'));
  expect(container).toMatchSnapshot();
});
