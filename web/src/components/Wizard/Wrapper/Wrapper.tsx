import {TWizardMap} from 'types/Wizard.types';
import WizardProvider from 'providers/Wizard/Wizard.provider';
import TracingBackend, {TracingBackendTab} from '../Steps/TracingBackend';
import CreateTest, {CreateTestTab} from '../Steps/CreateTest';

const steps: TWizardMap = {
  agent: undefined,
  tracing_backend: {
    name: 'Configure access to your OTel traces',
    description: '',
    component: TracingBackend,
    tabComponent: TracingBackendTab,
    isEnabled: false,
  },
  create_test: {
    name: 'Run your first test',
    description: '',
    component: CreateTest,
    tabComponent: CreateTestTab,
    isEnabled: false,
  },
};

interface IProps {
  children: React.ReactNode;
}

const Wrapper = ({children}: IProps) => <WizardProvider stepsMap={steps}>{children}</WizardProvider>;

export default Wrapper;
