import {Navigate, Route, Routes} from 'react-router-dom';
import {HistoryRouter} from 'redux-first-history/rr6';
import {history} from 'redux/store';

import Home from 'pages/Home';
import Test from 'pages/Test';
import RunDetail from 'pages/RunDetail';
import CreateTest from 'pages/CreateTest';

const {serverPathPrefix = '/'} = window.ENV || {};

const Router = () => (
  <HistoryRouter history={history} basename={serverPathPrefix}>
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/test/create" element={<CreateTest />} />
      <Route path="/test/:testId" element={<Test />} />

      <Route path="/test/:testId/run/:runId" element={<RunDetail />}>
        <Route path=":mode" element={<RunDetail />} />
      </Route>

      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  </HistoryRouter>
);

export default Router;
