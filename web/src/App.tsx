import * as Sentry from '@sentry/react';
import {useMemo} from 'react';
import {HistoryRouter} from 'redux-first-history/rr6';
import {useNavigate} from 'react-router-dom';

import ErrorBoundary from 'components/ErrorBoundary';
import {theme} from 'constants/Theme.constants';
import DashboardProvider from 'providers/Dashboard';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';
import {history} from 'redux/store';
import {ThemeProvider} from 'styled-components';
import Env from 'utils/Env';
import './App.css';
import BaseApp from './BaseApp';

const serverPathPrefix = Env.get('serverPathPrefix');

const App = () => {
  const navigate = useNavigate();
  const dashboardProviderValue = useMemo(() => ({baseUrl: '', navigate}), [navigate]);

  return (
    <ThemeProvider theme={theme}>
      <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
        <ReduxWrapperProvider>
          <HistoryRouter history={history} basename={serverPathPrefix}>
            <DashboardProvider value={dashboardProviderValue}>
              <BaseApp />
            </DashboardProvider>
          </HistoryRouter>
        </ReduxWrapperProvider>
      </Sentry.ErrorBoundary>
    </ThemeProvider>
  );
};

export default App;
