import {ICreateTestStep} from 'types/Plugins.types';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import PluginsComponentMap from './Plugins';

interface IProps {
  step: ICreateTestStep;
}

const CreateTestFactory = ({step: {component}}: IProps) => {
  const {pluginName} = useCreateTest();
  const Step = PluginsComponentMap[pluginName][component];

  return <Step />;
};

export default CreateTestFactory;
