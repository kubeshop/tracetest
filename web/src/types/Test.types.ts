import {FormInstance} from 'antd';
import {VariableDefinition, Request} from 'postman-collection';
import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {TriggerTypes} from 'constants/Test.constants';
import {HTTP_METHOD, SupportedPlugins} from 'constants/Common.constants';
import {Model, TGrpcSchemas, THttpSchemas, TTestSchemas, TTriggerSchemas} from './Common.types';
import {TTestSpecs} from './TestSpecs.types';

import {ICreateTestStep, IPlugin} from './Plugins.types';

export type TRequestAuth = THttpSchemas['HTTPRequest']['auth'];
export type TMethod = THttpSchemas['HTTPRequest']['method'];
export type TRawHeader = THttpSchemas['HTTPHeader'];
export type TRawGRPCHeader = TGrpcSchemas['GRPCHeader'];
export type TTriggerType = Required<TTriggerSchemas['Trigger']['triggerType']>;

export type TRawHTTPRequest = THttpSchemas['HTTPRequest'];
export type THeader = Model<TRawHeader, {}>;
export type THTTPRequest = Model<
  TRawHTTPRequest,
  {
    headers: THeader[];
  }
>;

export type TRawTriggerResult = TTriggerSchemas['TriggerResult'];
export type TTriggerResult = {
  type: TriggerTypes;
  headers?: THeader[];
  body?: string;
  statusCode: number;
  bodyMimeType?: string;
};

export type TRawGRPCRequest = TGrpcSchemas['GRPCRequest'];
export type TGRPCRequest = Model<
  TRawGRPCRequest,
  {
    metadata: Model<TRawGRPCHeader, {}>[];
  }
>;

export type TTriggerRequest = THTTPRequest | TRawGRPCRequest;

export type TRawTrigger = TTriggerSchemas['Trigger'];
export type TTrigger = {
  type: TriggerTypes;
  entryPoint: string;
  method: string;
  request: TTriggerRequest;
};

export type TRawTestSummary = TTestSchemas['TestSummary'];
export type TSummary = {
  runs: number;
  lastRun: {
    time: string;
    passes: number;
    fails: number;
  };
};

export type TRawTest = TTestSchemas['Test'];
export type TTest = Model<
  TRawTest,
  {
    definition: TTestSpecs;
    serviceUnderTest?: undefined;
    trigger: TTrigger;
    specs?: TTestSpecs;
    summary: TSummary;
  }
>;

export interface IRpcValues {
  message: string;
  auth: TRequestAuth;
  metadata: TGRPCRequest['metadata'];
  url: string;
  method: string;
  protoFile: File;
}

export interface IHttpValues {
  body: string;
  auth: TRequestAuth;
  headers: THTTPRequest['headers'];
  method: HTTP_METHOD;
  url: string;
}

export interface RequestDefinitionExtended extends Request {
  id: string;
  name: string;
}

export interface IPostmanValues extends IHttpValues {
  collectionFile?: File;
  envFile?: File;
  collectionTest?: string;
  requests: RequestDefinitionExtended[];
  variables: VariableDefinition[];
}

export interface ICurlValues extends IHttpValues {
  command: string;
}

export interface IBasicValues {
  name: string;
  description: string;
  testSuite: string;
}

export type TTestRequestDetailsValues = IRpcValues | IHttpValues | IPostmanValues | ICurlValues;
export type TDraftTest<T = TTestRequestDetailsValues> = Partial<IBasicValues & T>;
export type TDraftTestForm<T = TTestRequestDetailsValues> = FormInstance<TDraftTest<T>>;

export interface ITriggerService {
  getRequest(values: TDraftTest): Promise<TTriggerRequest>;
  validateDraft(draft: TDraftTest): Promise<boolean>;
  getInitialValues?(draft: TTriggerRequest): TDraftTest;
}

export interface ICreateTestState {
  draftTest: TDraftTest;
  stepList: ICreateTestStep[];
  stepNumber: number;
  pluginName: SupportedPlugins;
  isFormValid: boolean;
}

export type TCreateTestSliceActions = {
  reset: CaseReducer<ICreateTestState>;
  setPlugin: CaseReducer<ICreateTestState, PayloadAction<{plugin: IPlugin}>>;
  setStepNumber: CaseReducer<ICreateTestState, PayloadAction<{stepNumber: number; completeStep?: boolean}>>;
  setDraftTest: CaseReducer<ICreateTestState, PayloadAction<{draftTest: TDraftTest}>>;
  setIsFormValid: CaseReducer<ICreateTestState, PayloadAction<{isValid: boolean}>>;
};
