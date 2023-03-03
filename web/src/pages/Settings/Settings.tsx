import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import DataStoreProvider from 'providers/DataStore';
import SettingsProvider from 'providers/Settings';
import Content from './Content';

const Settings = () => (
  <Layout hasMenu>
    <DataStoreProvider>
      <SettingsProvider>
        <Content />
      </SettingsProvider>
    </DataStoreProvider>
  </Layout>
);

export default withAnalytics(Settings, 'settings');
