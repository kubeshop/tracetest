import {SupportedPlugins} from 'constants/Plugins.constants';
import {ICreateTestStep} from 'types/Plugins.types';
import PluginsComponentMap from '../CreateTestPlugins/Plugins';

interface IProps {
  step: ICreateTestStep;
  pluginName: SupportedPlugins;
}

const CreateTestStepFactory = ({step: {component}, pluginName}: IProps) => {
  const Step = PluginsComponentMap[pluginName][component];

  return <Step />;
};

export default CreateTestStepFactory;
