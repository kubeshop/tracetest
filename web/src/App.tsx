import * as Sentry from '@sentry/react';
import {ThemeProvider} from 'styled-components';

import ErrorBoundary from 'components/ErrorBoundary';
import GuidedTour from 'components/GuidedTour/GuidedTour';
import Router from 'components/Navigation';
import {theme} from 'constants/Theme.constants';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';

import './App.less';

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
        <GuidedTour>
          <ReduxWrapperProvider>
            <Router />
          </ReduxWrapperProvider>
        </GuidedTour>
      </Sentry.ErrorBoundary>
    </ThemeProvider>
  );
};

export default App;
