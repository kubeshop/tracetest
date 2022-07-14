import {BrowserRouter, Routes, Route, Navigate} from 'react-router-dom';
import Home from 'pages/Home';
import Test from 'pages/Test';
import Trace from 'pages/Trace';
import CreateTest from 'pages/CreateTest';
import EditTest from 'pages/EditTest';

const {serverPathPrefix = '/'} = window.ENV || {};

const Router = (): JSX.Element => {
  return (
    <BrowserRouter basename={serverPathPrefix}>
      <Routes>
        <Route path="/" element={<Home path="home" />} />
        <Route path="/test/:testId/edit" element={<EditTest path="edit-test" />} />
        <Route path="/test/create" element={<CreateTest path="create-test" />} />
        <Route path="/test/:testId" element={<Test path="test-details" />} />
        <Route path="/test/:testId/run/:runId" element={<Trace path="trace" />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
