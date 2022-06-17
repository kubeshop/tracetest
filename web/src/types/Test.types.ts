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
