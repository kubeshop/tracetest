import {Navigate, Route, Routes} from 'react-router-dom';

import Layout from 'components/Layout/Layout';
import VariableSet from 'pages/VariableSet';
import Home from 'pages/Home';
import TestSuites from 'pages/TestSuites';
import RunDetail from 'pages/RunDetail';
import Settings from 'pages/Settings';
import Test from 'pages/Test';
import TestSuite from 'pages/TestSuite';
import TestSuiteRunOverview from 'pages/TestSuiteRunOverview';
import TestSuiteRunAutomate from 'pages/TestSuiteRunAutomate';
import AutomatedTestRun from 'pages/AutomatedTestRun';
import CreateTest from 'pages/CreateTest';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

const Router = () => {
  const {baseUrl} = useDashboard();

  return (
    <Routes>
      <Route element={<Layout hasMenu />}>
        <Route path="/" element={<Home />} />
        <Route path="/testsuites" element={<TestSuites />} />
        <Route path="/variablesets" element={<VariableSet />} />
        <Route path="/settings" element={<Settings />} />
      </Route>

      <Route element={<Layout />}>
        <Route path="/test/create/:triggerType" element={<CreateTest />} />
        <Route path="/test/:testId" element={<Test />} />
        <Route path="/test/:testId/run/:runId" element={<RunDetail />} />
        <Route path="/test/:testId/run/:runId/:mode" element={<RunDetail />} />
        <Route path="/test/:testId/run" element={<AutomatedTestRun />} />

        <Route path="/testsuite/:testSuiteId" element={<TestSuite />} />
        <Route path="/testsuite/:testSuiteId/run/:runId" element={<TestSuiteRunOverview />} />
        <Route path="/testsuite/:testSuiteId/run/:runId/overview" element={<TestSuiteRunOverview />} />
        <Route path="/testsuite/:testSuiteId/run/:runId/automate" element={<TestSuiteRunAutomate />} />
      </Route>

      <Route path="*" element={<Navigate to={baseUrl} />} />
    </Routes>
  );
};

export default Router;
