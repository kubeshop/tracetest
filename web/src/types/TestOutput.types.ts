import {Model, TTestSchemas} from './Common.types';

export type TRawTestOutput = TTestSchemas['TestOutput'];

export type TTestOutput = {
  isDeleted: boolean;
  isDraft: boolean;
  name: string;
  selector: string;
  value: string;
  valueRun: string;
  valueRunDraft: string;
  id: number;
  spanId: string;
};

export type TRawTestRunOutput = {
  name?: string;
  value?: string;
  spanId?: string;
};

export type TTestRunOutput = Model<TRawTestRunOutput, {}>;
