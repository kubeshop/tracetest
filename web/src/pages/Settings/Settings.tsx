import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import DataStoreProvider from 'providers/DataStore';
import Content from './Content';

const Settings = () => (
  <Layout hasMenu>
    <DataStoreProvider>
      <Content />
    </DataStoreProvider>
  </Layout>
);

export default withAnalytics(Settings, 'settings');
