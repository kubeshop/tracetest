import Envs from 'pages/Environments';

import Home from 'pages/Home';
import RunDetail from 'pages/RunDetail';
import Test from 'pages/Test';
import {Navigate, Route, Routes} from 'react-router-dom';
import {HistoryRouter} from 'redux-first-history/rr6';
import {history} from 'redux/store';
import ExperimentalFeature from 'utils/ExperimentalFeature';

const {serverPathPrefix = '/'} = window.ENV || {};
const isTransctionFeatureEnabled = ExperimentalFeature.isEnabled('transactions');

const Router = () => (
  <HistoryRouter history={history} basename={serverPathPrefix}>
    <Routes>
      <Route path="/" element={<Home />} />
      {isTransctionFeatureEnabled && <Route path="/environments" element={<Envs />} />}
      <Route path="/test/:testId" element={<Test />} />

      <Route path="/test/:testId/run/:runId" element={<RunDetail />}>
        <Route path=":mode" element={<RunDetail />} />
      </Route>

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  </HistoryRouter>
);

export default Router;
