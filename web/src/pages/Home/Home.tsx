import Layout from 'components/Layout';
import HomeContent from './HomeContent';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';

const Home = (): JSX.Element => {
  return (
    <Layout>
      <HomeContent />
    </Layout>
  );
};

export default withAnalytics(Home, 'home');
