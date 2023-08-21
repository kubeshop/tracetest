import Router from 'components/Router';
import SettingsValuesProvider from 'providers/SettingsValues';
import {TCustomHeader} from 'components/Layout/Layout';

interface IProps {
  customHeader?: TCustomHeader;
}

const BaseApp = ({customHeader}: IProps) => (
  <SettingsValuesProvider>
    <Router customHeader={customHeader} />
  </SettingsValuesProvider>
);

export default BaseApp;
