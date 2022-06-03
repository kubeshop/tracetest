import * as Sentry from '@sentry/react';
import './App.less';
import AnalyticsProvider from './components/Analytics/AnalyticsProvider';
import ErrorBoundary from './components/ErrorBoundary';
import {GuidedTourProvider} from './components/GuidedTour/GuidedTourProvider';
import Router from './components/Navigation';
import {ReduxWrapperProvider} from './redux/ReduxWrapperProvider';

const App = () => {
  return (
    <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
      <GuidedTourProvider>
        <AnalyticsProvider>
          <ReduxWrapperProvider>
            <Router />
          </ReduxWrapperProvider>
        </AnalyticsProvider>
      </GuidedTourProvider>
    </Sentry.ErrorBoundary>
  );
};

export default App;
