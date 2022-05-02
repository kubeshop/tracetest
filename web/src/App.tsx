import {Provider} from 'react-redux';
import * as Sentry from '@sentry/react';
import {TourProvider} from '@reactour/tour';
import Router from './components/Navigation';
import {store} from './redux/store';
import './App.css';
import AnalyticsProvider from './components/Analytics/AnalyticsProvider';
import ErrorBoundary from './components/ErrorBoundary';

const App = () => {
  return (
    <TourProvider steps={[]}>
      <AnalyticsProvider>
        <Provider store={store}>
          <Router />
        </Provider>
      </AnalyticsProvider>
    </TourProvider>
  );
};

export default Sentry.withErrorBoundary(App, {
  fallback: <ErrorBoundary />,
});
