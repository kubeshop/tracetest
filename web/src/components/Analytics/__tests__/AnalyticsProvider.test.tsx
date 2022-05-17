import {render} from '@testing-library/react';
import {Provider} from 'react-redux';
import {store} from '../../../redux/store';
import AnalyticsProvider from '../index';

test('AnalyticsProvider', () => {
  const {container, getAllByTestId} = render(
    <Provider store={store}>
      <AnalyticsProvider>
        <h2 data-cy="sample">Cesco</h2>
      </AnalyticsProvider>
    </Provider>
  );
  expect(container).toMatchSnapshot();
  expect(getAllByTestId('not_initialized').length).toBe(1);
});
