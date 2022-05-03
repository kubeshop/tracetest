import {Modify, Schemas} from './Common.types';
import {TRawTrace, TTrace} from './Trace.types';

export type TRawTestRunResult = Modify<
  Schemas['TestRunResult'],
  {
    trace: TRawTrace;
  }
>;

export type TTestRunResult = Modify<TRawTestRunResult, {trace?: TTrace}>;
export type TTestState = Schemas['TestRunResult']['state'];
