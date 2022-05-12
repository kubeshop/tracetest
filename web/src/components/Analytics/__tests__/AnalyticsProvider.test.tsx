import {render} from '@testing-library/react';
import {Provider} from 'react-redux';
import {store} from '../../../redux/store';
import AnalyticsProvider from '../index';

test('AnalyticsProvider', () => {
  const result = render(
    <Provider store={store}>
      <AnalyticsProvider>
        <h2>Cesco</h2>
      </AnalyticsProvider>
    </Provider>
  );
  expect(result.container).toMatchSnapshot();
});
