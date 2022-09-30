import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import CreateTransactionProvider from 'providers/CreateTransaction';
import CreateTestProvider from 'providers/CreateTest';
import ExperimentalFeature from 'utils/ExperimentalFeature';
import Content from './Content';

const Home = () => (
  <Layout hasMenu={ExperimentalFeature.isEnabled('transactions')}>
    <CreateTransactionProvider>
      <CreateTestProvider>
        <Content />
      </CreateTestProvider>
    </CreateTransactionProvider>
  </Layout>
);

export default withAnalytics(Home, 'home');
