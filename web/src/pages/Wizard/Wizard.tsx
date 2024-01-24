import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import DataStoreProvider from 'providers/DataStore/DataStore.provider';
import SettingsProvider from 'providers/Settings/Settings.provider';
import HomeContent from './HomeContent';

const Wizard = () => (
  <DataStoreProvider>
    <SettingsProvider>
      <HomeContent />
    </SettingsProvider>
  </DataStoreProvider>
);

export default withAnalytics(Wizard, 'wizard');
