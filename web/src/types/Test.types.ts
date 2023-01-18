import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {BaseQueryFn, FetchArgs, FetchBaseQueryError, FetchBaseQueryMeta} from '@reduxjs/toolkit/dist/query';
import {EndpointBuilder} from '@reduxjs/toolkit/dist/query/endpointDefinitions';
import {FormInstance} from 'antd';
import {VariableDefinition, Request} from 'postman-collection';

import {HTTP_METHOD, SupportedPlugins} from 'constants/Common.constants';
import {TracetestApiTags, TriggerTypes} from 'constants/Test.constants';
import {Model, TGrpcSchemas, THttpSchemas, TTestSchemas, TTraceIDSchemas, TTriggerSchemas} from './Common.types';
import {ICreateTestStep, IPlugin} from './Plugins.types';
import {TTestOutput} from './TestOutput.types';
import {TTestSpecs} from './TestSpecs.types';

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

export type TRawTRACEIDRequest = TTraceIDSchemas['TRACEIDRequest'];

export type TTRACEIDRequest = Model<
  TRawTRACEIDRequest,
  {
    id: string;
  }
>;

export type TTriggerRequest = THTTPRequest | TRawGRPCRequest | TRawTRACEIDRequest;

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
    outputs?: TTestOutput[];
    createdAt?: string;
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

export interface ITraceIDValues extends IHttpValues {
  id: string;
}

export type TTestRequestDetailsValues = IRpcValues | IHttpValues | IPostmanValues | ICurlValues | ITraceIDValues;
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

export type TTestApiEndpointBuilder = EndpointBuilder<
  BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError, {}, FetchBaseQueryMeta>,
  TracetestApiTags,
  'tests'
>;
