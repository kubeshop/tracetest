import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import Content from './Content';

const Home = () => (
  <Layout hasMenu={ExperimentalFeature.isEnabled('transactions')}>
    <Content />
  </Layout>
);

export default withAnalytics(Home, 'home');
