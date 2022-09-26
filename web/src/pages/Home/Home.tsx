import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const {featuresEnabled = '[]'} = window.ENV || {};
const parsedFeaturesEnabled = JSON.parse(featuresEnabled);
const isTransactionsEnabled = parsedFeaturesEnabled?.includes('transactions');

const Home = () => (
  <Layout hasMenu={isTransactionsEnabled}>
    <Content />
  </Layout>
);

export default withAnalytics(Home, 'home');
