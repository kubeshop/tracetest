import {Navigate, Route, Routes} from 'react-router-dom';

import VariableSet from 'pages/VariableSet';
import Home from 'pages/Home';
import RunDetail from 'pages/RunDetail';
import Settings from 'pages/Settings';
import Test from 'pages/Test';
import Transaction from 'pages/Transaction';
import TransactionRunOverview from 'pages/TransactionRunOverview';
import TransactionRunAutomate from 'pages/TransactionRunAutomate';
import AutomatedTestRun from 'pages/AutomatedTestRun';

const Router = () => (
  <Routes>
    <Route path="/" element={<Home />} />

    <Route path="/variablesets" element={<VariableSet />} />

    <Route path="/settings" element={<Settings />} />

    <Route path="/test/:testId" element={<Test />} />
    <Route path="/test/:testId/run/:runId" element={<RunDetail />} />
    <Route path="/test/:testId/run/:runId/:mode" element={<RunDetail />} />
    <Route path="/test/:testId/run" element={<AutomatedTestRun />} />

    <Route path="/transaction/:transactionId" element={<Transaction />} />
    <Route path="/transaction/:transactionId/run/:runId" element={<TransactionRunOverview />} />
    <Route path="/transaction/:transactionId/run/:runId/overview" element={<TransactionRunOverview />} />
    <Route path="/transaction/:transactionId/run/:runId/automate" element={<TransactionRunAutomate />} />

    <Route path="*" element={<Navigate to="" />} />
  </Routes>
);

export default Router;
