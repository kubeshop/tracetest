import {useMemo} from 'react';
import {IWizardStep} from 'types/Wizard.types';
import WizardProvider from 'providers/Wizard/Wizard.provider';
import DataStoreProvider from 'providers/DataStore/DataStore.provider';
import SettingsProvider from 'providers/Settings/Settings.provider';
import TracingBackend from '../Steps/TracingBackend/TracingBackend';
import RunTest from '../Steps/RunTest';

interface IProps {
  children: React.ReactNode;
}

const Wrapper = ({children}: IProps) => {
  const steps = useMemo<IWizardStep[]>(
    () => [
      {
        id: 'tracing-backend',
        name: 'Setup your Tracing Backend',
        description: '',
        component: TracingBackend,
        status: 'pending', // grab status from somewhere else
      },
      {
        id: 'run-test',
        name: 'Run your first test',
        description: '',
        component: RunTest,
        status: 'pending',
      },
    ],
    []
  );

  return (
    <DataStoreProvider>
      <SettingsProvider>
        <WizardProvider steps={steps}>{children}</WizardProvider>
      </SettingsProvider>
    </DataStoreProvider>
  );
};

export default Wrapper;
