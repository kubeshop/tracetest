import {SupportedPlugins} from 'constants/Plugins.constants';
import {TriggerTypes} from '../constants/Test.constants';

export type TStepStatus = 'complete' | 'pending' | 'selected';

export interface ICreateTestStep {
  id: string;
  name: string;
  title: string;
  component: string;
  status?: TStepStatus;
  isDefaultValid?: boolean;
}

export interface IPlugin {
  name: SupportedPlugins;
  title: string;
  description: string;
  stepList: ICreateTestStep[];
  isActive: boolean;
  type: TriggerTypes;
}

export interface IPluginStepProps {}

export interface IPluginComponentMap extends Record<string, (props: IPluginStepProps) => React.ReactElement> {}
