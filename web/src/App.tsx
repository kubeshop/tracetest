import {Provider} from 'react-redux';
import {TourProvider} from '@reactour/tour';
import Router from './navigation';
import {store} from './redux/store';
import './App.css';
import AnalyticsProvider from './components/Analytics/AnalyticsProvider';

const App = () => {
  return (
    <TourProvider steps={[]}>
      <AnalyticsProvider>
        <Provider store={store}>
          <Router />
        </Provider>
      </AnalyticsProvider>
    </TourProvider>
  );
};

export default App;
