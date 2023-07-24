import Router from 'components/Router';
import SettingsValuesProvider from 'providers/SettingsValues';

const BaseApp = () => (
  <SettingsValuesProvider>
    <Router />
  </SettingsValuesProvider>
);

export default BaseApp;
