import {TWizardMap} from 'types/Wizard.types';
import WizardProvider from 'providers/Wizard/Wizard.provider';
import {withCustomization} from 'providers/Customization';
import DataStoreProvider from 'providers/DataStore/DataStore.provider';
import SettingsProvider from 'providers/Settings/Settings.provider';
import TracingBackend from '../Steps/TracingBackend';
import RunTest from '../Steps/RunTest';

const steps: TWizardMap = {
  agent: {
    name: '',
    description: '',
    component: () => <div>Agent</div>,
  },
  tracing_backend: {
    name: 'Configure access to your OTel traces',
    description: '',
    component: TracingBackend,
  },
  create_test: {
    name: 'Run your first test',
    description: '',
    component: RunTest,
  },
};

interface IProps {
  children: React.ReactNode;
}

const Wrapper = ({children}: IProps) => (
  <DataStoreProvider>
    <SettingsProvider>
      <WizardProvider stepsMap={steps}>{children}</WizardProvider>
    </SettingsProvider>
  </DataStoreProvider>
);

export default withCustomization(Wrapper, 'wizardWrapper');
