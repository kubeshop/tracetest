import {ICreateTestStep} from 'types/Plugins.types';
import TestSuitedPlugin from '.';

interface IProps {
  step: ICreateTestStep;
}

const CreateTestSuiteFactory = ({step: {component}}: IProps) => {
  const Step = TestSuitedPlugin[component];

  return <Step />;
};

export default CreateTestSuiteFactory;
