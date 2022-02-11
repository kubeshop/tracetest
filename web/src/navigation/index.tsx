import {BrowserRouter, Routes, Route} from 'react-router-dom';
import Home from '../pages/Home';
import Trace from '../pages/Trace';

const Router = (): JSX.Element => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/trace" element={<Trace />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
