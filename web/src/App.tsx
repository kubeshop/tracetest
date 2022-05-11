import {TourProvider} from '@reactour/tour';
import * as Sentry from '@sentry/react';
import './App.less';
import AnalyticsProvider from './components/Analytics/AnalyticsProvider';
import ErrorBoundary from './components/ErrorBoundary';
import Router from './components/Navigation';
import {ReduxWrapperProvider} from './redux/ReduxWrapperProvider';

const App = () => {
  return (
    <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
      <TourProvider steps={[]}>
        <AnalyticsProvider>
          <ReduxWrapperProvider>
            <Router />
          </ReduxWrapperProvider>
        </AnalyticsProvider>
      </TourProvider>
    </Sentry.ErrorBoundary>
  );
};

export default App;
