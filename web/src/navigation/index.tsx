import {BrowserRouter, Routes, Route} from 'react-router-dom';
import Home from '../pages/Home';
import Test from '../pages/Test';

const Router = (): JSX.Element => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/test/:id" element={<Test />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
