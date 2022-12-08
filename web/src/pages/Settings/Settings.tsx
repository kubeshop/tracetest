import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import SetupConfigProvider from 'providers/SetupConfig';
import Content from './Content';

const Settings = () => (
  <Layout hasMenu>
    <SetupConfigProvider>
      <Content />
    </SetupConfigProvider>
  </Layout>
);

export default withAnalytics(Settings, 'settings');
