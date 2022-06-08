import * as Sentry from '@sentry/react';
import AnalyticsProvider from 'components/Analytics/AnalyticsProvider';
import ErrorBoundary from 'components/ErrorBoundary';
import GuidedTour from 'components/GuidedTour/GuidedTour';
import Router from 'components/Navigation';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';
import './App.less';

const App = () => {
  return (
    <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
      <GuidedTour>
        <AnalyticsProvider>
          <ReduxWrapperProvider>
            <Router />
          </ReduxWrapperProvider>
        </AnalyticsProvider>
      </GuidedTour>
    </Sentry.ErrorBoundary>
  );
};

export default App;
