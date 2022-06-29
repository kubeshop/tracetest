import {Model, THttpSchemas, TTestSchemas} from './Common.types';
import {TTestDefinition} from './TestDefinition.types';

export type TRawTest = TTestSchemas['Test'];
export type TTest = Model<
  TRawTest,
  {
    definition: TTestDefinition;
    serviceUnderTest?: {
      request?: THttpSchemas['HTTPRequest'];
    };
  }
>;

export type TRequest = THttpSchemas['HTTPRequest'];
export type TRequestAuth = THttpSchemas['HTTPRequest']['auth'];
export type TMethod = THttpSchemas['HTTPRequest']['method'];
