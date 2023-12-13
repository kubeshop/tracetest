import {SupportedPlugins} from 'constants/Common.constants';
import {TriggerTypes} from '../constants/Test.constants';
import {TDraftTest} from './Test.types';

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
  demoList: TDraftTest[];
  isActive: boolean;
  type: TriggerTypes;
  requestType?: TriggerTypes;
}

export interface IPluginStepProps {}

export interface IPluginComponentMap extends Record<string, (props: IPluginStepProps) => React.ReactElement> {}
