import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const Home = () => (
  <Layout hasMenu>
    <Content />
  </Layout>
);

export default withAnalytics(Home, 'home');
