import {BrowserRouter, Routes, Route} from 'react-router-dom';
import Home from '../pages/Home';
import Trace from '../pages/Trace';
import Test from '../pages/Test';

const Router = (): JSX.Element => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/trace" element={<Trace />} />
        <Route path="/test" element={<Test />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
