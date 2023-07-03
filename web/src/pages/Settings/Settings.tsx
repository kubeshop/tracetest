import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import DataStoreProvider from 'providers/DataStore';
import SettingsProvider from 'providers/Settings';
import Content from './Content';
import ContactUs from '../../components/ContactUs/ContactUs';

const Settings = () => (
  <Layout hasMenu>
    <ContactUs>
      <DataStoreProvider>
        <SettingsProvider>
          <Content />
        </SettingsProvider>
      </DataStoreProvider>
    </ContactUs>
  </Layout>
);

export default withAnalytics(Settings, 'settings');
