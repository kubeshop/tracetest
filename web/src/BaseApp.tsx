import Router from 'components/Router';
import SettingsValuesProvider from 'providers/SettingsValues';
import CreateTestProvider from 'providers/CreateTest/CreateTest.provider';

const BaseApp = () => (
  <CreateTestProvider>
    <SettingsValuesProvider>
      <Router />
    </SettingsValuesProvider>
  </CreateTestProvider>
);

export default BaseApp;
