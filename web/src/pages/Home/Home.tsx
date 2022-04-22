import Layout from 'components/Layout';
import {TourProvider} from '@reactour/tour';
import HomeContent from './HomeContent';

const Home = (): JSX.Element => {
  return (
    <Layout>
      <TourProvider steps={[]}>
        <HomeContent />
      </TourProvider>
    </Layout>
  );
};

export default Home;
