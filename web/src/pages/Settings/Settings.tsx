import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import DataStoreProvider from 'providers/DataStore';
import SettingsProvider from 'providers/Settings';
import Content from './Content';
import ContactUs from '../../components/ContactUs/ContactUs';

const Settings = () => (
  <ContactUs>
    <DataStoreProvider>
      <SettingsProvider>
        <Content />
      </SettingsProvider>
    </DataStoreProvider>
  </ContactUs>
);

export default withAnalytics(Settings, 'settings');
