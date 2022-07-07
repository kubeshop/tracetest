import {TriggerTypes} from 'constants/Test.constants';
import {Model, TGrpcSchemas, THttpSchemas, TTestSchemas, TTriggerSchemas} from './Common.types';
import {TTestDefinition} from './TestDefinition.types';

export type TRequestAuth = THttpSchemas['HTTPRequest']['auth'];
export type TMethod = THttpSchemas['HTTPRequest']['method'];
export type TRawHeader = THttpSchemas['HTTPHeader'];
export type TRawGRPCHeader = TGrpcSchemas['GRPCHeader'];
export type TTriggerType = Required<TTriggerSchemas['Trigger']['triggerType']>;

export type TRawHTTPRequest = THttpSchemas['HTTPRequest'];
export type THTTPRequest = Model<
  TRawHTTPRequest,
  {
    headers: Model<TRawHeader, {}>[];
  }
>;

export type TRawGRPCRequest = TGrpcSchemas['GRPCRequest'];
export type TGRPCRequest = Model<
  TRawGRPCRequest,
  {
    metadata: Model<TRawGRPCHeader, {}>[];
  }
>;

export type TRawTrigger = TTriggerSchemas['Trigger'];
export type TTrigger = {
  type: TriggerTypes;
  entryPoint: string;
  method: string;
  request: THTTPRequest | TRawGRPCRequest;
};

export type TRawTest = TTestSchemas['Test'];
export type TTest = Model<
  TRawTest,
  {
    definition: TTestDefinition;
    serviceUnderTest?: undefined;
    trigger: TTrigger;
  }
>;
