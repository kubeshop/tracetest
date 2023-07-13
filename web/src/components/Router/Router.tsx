import {Navigate, Route, Routes} from 'react-router-dom';
import {HistoryRouter} from 'redux-first-history/rr6';
import {history} from 'redux/store';

import Envs from 'pages/Environments';
import Home from 'pages/Home';
import RunDetail from 'pages/RunDetail';
import Settings from 'pages/Settings';
import Test from 'pages/Test';
import Transaction from 'pages/Transaction';
import TransactionRunTrigger from 'pages/TransactionRunTrigger';
import TransactionRunAutomate from 'pages/TransactionRunAutomate';
import AutomatedTestRun from 'pages/AutomatedTestRun';
import Env from 'utils/Env';

const serverPathPrefix = Env.get('serverPathPrefix');

const Router = () => (
  <HistoryRouter history={history} basename={serverPathPrefix}>
    <Routes>
      <Route path="/" element={<Home />} />

      <Route path="/environments" element={<Envs />} />

      <Route path="/settings" element={<Settings />} />

      <Route path="/test/:testId" element={<Test />} />
      <Route path="/test/:testId/run/:runId" element={<RunDetail />} />
      <Route path="/test/:testId/run/:runId/:mode" element={<RunDetail />} />
      <Route path="/test/:testId/run" element={<AutomatedTestRun />} />

      <Route path="/transaction/:transactionId" element={<Transaction />} />
      <Route path="/transaction/:transactionId/run/:runId" element={<TransactionRunTrigger />} />
      <Route path="/transaction/:transactionId/run/:runId/trigger" element={<TransactionRunTrigger />} />
      <Route path="/transaction/:transactionId/run/:runId/automate" element={<TransactionRunAutomate />} />

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  </HistoryRouter>
);

export default Router;
