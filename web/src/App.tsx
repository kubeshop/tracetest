import * as Sentry from '@sentry/react';
import {HistoryRouter} from 'redux-first-history/rr6';

import CustomizationWrapper from 'components/CustomizationWrapper';
import DashboardWrapper from 'components/DashboardWrapper';
import ErrorBoundary from 'components/ErrorBoundary';
import {theme} from 'constants/Theme.constants';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';
import {history} from 'redux/store';
import {ThemeProvider} from 'styled-components';
import Env from 'utils/Env';
import './App.css';
import BaseApp from './BaseApp';

const serverPathPrefix = Env.get('serverPathPrefix');

const App = () => (
  <ThemeProvider theme={theme}>
    <Sentry.ErrorBoundary fallback={({error}) => <ErrorBoundary error={error} />}>
      <ReduxWrapperProvider>
        <HistoryRouter history={history} basename={serverPathPrefix}>
          <CustomizationWrapper>
            <DashboardWrapper>
              <BaseApp />
            </DashboardWrapper>
          </CustomizationWrapper>
        </HistoryRouter>
      </ReduxWrapperProvider>
    </Sentry.ErrorBoundary>
  </ThemeProvider>
);

export default App;
