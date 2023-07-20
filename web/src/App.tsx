import * as Sentry from '@sentry/react';

import ErrorBoundary from 'components/ErrorBoundary';
import Router from 'components/Router';
import {theme} from 'constants/Theme.constants';
import SettingsValuesProvider from 'providers/SettingsValues';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';
import {ThemeProvider} from 'styled-components';
import './App.css';

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
        <ReduxWrapperProvider>
          <SettingsValuesProvider>
            <Router />
          </SettingsValuesProvider>
        </ReduxWrapperProvider>
      </Sentry.ErrorBoundary>
    </ThemeProvider>
  );
};

export default App;
