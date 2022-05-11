import {BrowserRouter, Routes, Route, Navigate} from 'react-router-dom';
import Home from 'pages/Home';
import Test from 'pages/Test';
import Trace from '../../pages/Trace';

const Router = (): JSX.Element => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home path="home" />} />
        <Route path="/test/:id" element={<Test path="test-details" />} />
        <Route path="/test/:testId/result/:resultId" element={<Trace path="trace" />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
