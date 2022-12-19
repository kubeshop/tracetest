import * as Sentry from '@sentry/react';

import ErrorBoundary from 'components/ErrorBoundary';
import Router from 'components/Router';
import {theme} from 'constants/Theme.constants';
import DataStoreConfigProvider from 'providers/DataStoreConfig';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';
import {ThemeProvider} from 'styled-components';
import './App.less';

const App = () => {
  return (
    <ThemeProvider theme={theme}>
      <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
        <ReduxWrapperProvider>
          <DataStoreConfigProvider>
            <Router />
          </DataStoreConfigProvider>
        </ReduxWrapperProvider>
      </Sentry.ErrorBoundary>
    </ThemeProvider>
  );
};

export default App;
