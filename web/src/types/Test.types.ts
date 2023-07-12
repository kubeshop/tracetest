import {CaseReducer, PayloadAction} from '@reduxjs/toolkit';
import {BaseQueryFn, FetchArgs, FetchBaseQueryError, FetchBaseQueryMeta} from '@reduxjs/toolkit/dist/query';
import {EndpointBuilder} from '@reduxjs/toolkit/dist/query/endpointDefinitions';
import {FormInstance} from 'antd';
import {VariableDefinition, Request} from 'postman-collection';

import {HTTP_METHOD, SupportedPlugins} from 'constants/Common.constants';
import {TracetestApiTags} from 'constants/Test.constants';
import {Model, TGrpcSchemas, THttpSchemas} from './Common.types';
import {ICreateTestStep, IPlugin} from './Plugins.types';
import GRPCRequest from '../models/GrpcRequest.model';
import HttpRequest from '../models/HttpRequest.model';
import TraceIDRequest from '../models/TraceIDRequest.model';

export type TRequestAuth = THttpSchemas['HTTPRequest']['auth'];
export type TMethod = THttpSchemas['HTTPRequest']['method'];
export type TRawHeader = THttpSchemas['HTTPHeader'];
export type TRawGRPCHeader = TGrpcSchemas['GRPCHeader'];
export type THeader = Model<TRawHeader, {}>;

export interface IRpcValues {
  message: string;
  auth: TRequestAuth;
  metadata: GRPCRequest['metadata'];
  url: string;
  method: string;
  protoFile: File;
}

export interface IHttpValues {
  body: string;
  auth: TRequestAuth;
  headers: HttpRequest['headers'];
  method: HTTP_METHOD;
  url: string;
  sslVerification: boolean;
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

export type TTriggerRequest = HttpRequest | GRPCRequest | TraceIDRequest;
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
