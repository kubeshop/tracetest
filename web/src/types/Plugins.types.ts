import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {SupportedPlugins} from 'constants/Plugins.constants';
import {TRecursivePartial} from './Common.types';
import {TRawTest, TTriggerType} from './Test.types';

export type TStepStatus = 'complete' | 'pending' | 'selected';
export type TDraftTest = TRecursivePartial<TRawTest>;
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
  type: TTriggerType;
}

export interface IPluginStepProps {}

export interface IPluginComponentMap extends Record<string, (props: IPluginStepProps) => React.ReactElement> {}

export interface ICreateTestState {
  draftTest: TDraftTest;
  stepList: ICreateTestStep[];
  stepNumber: number;
  pluginName: SupportedPlugins;
}

export type TCreateTestSliceActions = {
  setPlugin: CaseReducer<ICreateTestState, PayloadAction<{plugin: IPlugin}>>;
  setStepNumber: CaseReducer<ICreateTestState, PayloadAction<{stepNumber: number; completeStep?: boolean}>>;
  setDraftTest: CaseReducer<ICreateTestState, PayloadAction<{draftTest: TDraftTest}>>;
};
