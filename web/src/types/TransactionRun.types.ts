import {Model, TTransactionsSchemas} from './Common.types';
import {TEnvironment} from './Environment.types';
import {TTest} from './Test.types';
import {TTestRun} from './TestRun.types';

export type TRawTransactionRun = TTransactionsSchemas['TransactionRun'];

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
