import {render} from '@testing-library/react';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import AnalyticsProvider from '../../Analytics';
import Router from '../Router';

test('Router', async () => {
  const {container} = render(
    <AnalyticsProvider>
      <ReduxWrapperProvider>
        <Router />
      </ReduxWrapperProvider>
    </AnalyticsProvider>
  );

  expect(container).toMatchSnapshot();
});
