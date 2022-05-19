import {render, waitFor} from '@testing-library/react';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import Router from '../Router';

jest.mock('../../../services/Analytics/Analytics.service', () => {
  return {
    event: jest.fn(),
  };
});

test('Router', async () => {
  const {container, getByTestId} = render(
    <ReduxWrapperProvider>
      <Router />
    </ReduxWrapperProvider>
  );

  await waitFor(() => getByTestId('github-link'));
  expect(container).toMatchSnapshot();
});
