import {Model} from './Common.types';
import {TEnvironment, TRawEnvironment} from './Environment.types';
import {TRawTest, TTest} from './Test.types';
import {TRawTestRun, TTestRun} from './TestRun.types';

export type TRawTransactionRun = {
  id: string;
  createdAt: string;
  completedAt: string;
  state: 'CREATED' | 'EXECUTING' | 'FINISHED' | 'FAILED';
  steps: TRawTest[];
  stepRuns: TRawTestRun[];
  environment?: TRawEnvironment;
  metadata?: {[key: string]: string};
};

export type TTransactionRun = Model<
  TRawTransactionRun,
  {
    steps: TTest[];
    stepRuns: TTestRun[];
    environment?: TEnvironment;
    metadata?: {[key: string]: string};
  }
>;

export type TTransactionRunState = TRawTransactionRun['state'];
