import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import CreateTransactionProvider from 'providers/CreateTransaction';
import CreateTestProvider from 'providers/CreateTest';
import Content from './Content';

const Home = () => (
  <Layout hasMenu>
    <CreateTransactionProvider>
      <CreateTestProvider>
        <Content />
      </CreateTestProvider>
    </CreateTransactionProvider>
  </Layout>
);

export default withAnalytics(Home, 'home');
