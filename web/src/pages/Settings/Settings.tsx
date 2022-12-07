import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const Settings = () => (
  <Layout hasMenu>
    <Content />
  </Layout>
);

export default withAnalytics(Settings, 'settings');
