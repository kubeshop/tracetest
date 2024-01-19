import Router from 'components/Router';
import SettingsValuesProvider from 'providers/SettingsValues';
import Wrapper from './components/Wizard/Wrapper';

const BaseApp = () => (
  <Wrapper>
    <SettingsValuesProvider>
      <Router />
    </SettingsValuesProvider>
  </Wrapper>
);

export default BaseApp;
