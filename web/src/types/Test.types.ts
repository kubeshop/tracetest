import {FormInstance} from 'antd';
import {VariableDefinition, Request} from 'postman-collection';

import {HTTP_METHOD} from 'constants/Common.constants';
import GRPCRequest from 'models/GrpcRequest.model';
import HttpRequest from 'models/HttpRequest.model';
import TraceIDRequest from 'models/TraceIDRequest.model';
import KafkaRequest from 'models/KafkaRequest.model';
import {Model, TGrpcSchemas, THttpSchemas, TKafkaSchemas} from './Common.types';

export type TRequestAuth = THttpSchemas['HTTPRequest']['auth'];
export type TMethod = THttpSchemas['HTTPRequest']['method'];
export type TRawHeader = THttpSchemas['HTTPHeader'];
export type TRawGRPCHeader = TGrpcSchemas['GRPCHeader'];
export type THeader = Model<TRawHeader, {}>;

export type TKafkaRequestAuth = TKafkaSchemas['KafkaRequest']['authentication'];
export type TRawKafkaMessageHeader = TKafkaSchemas['KafkaMessageHeader'];

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

export interface IKafkaValues {
  brokerUrls: string[];
  topic: string;
  authentication: TKafkaRequestAuth;
  sslVerification: boolean;
  headers: KafkaRequest['headers'];
  messageKey: string;
  messageValue: string;
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

export type TTriggerRequest = HttpRequest | GRPCRequest | TraceIDRequest | KafkaRequest;
export interface ITriggerService {
  getRequest(values: TDraftTest): Promise<TTriggerRequest>;
  validateDraft(draft: TDraftTest): Promise<boolean>;
  getInitialValues?(draft: TTriggerRequest): TDraftTest;
}

export interface IImportService {
  getRequest(values: TDraftTest): Promise<TDraftTest>;
  validateDraft(draft: TDraftTest): Promise<boolean>;
}
