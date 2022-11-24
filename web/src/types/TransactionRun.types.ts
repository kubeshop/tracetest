import {Model, TTransactionsSchemas} from './Common.types';
import {TEnvironment} from './Environment.types';
import {TTestRun} from './TestRun.types';

export type TRawTransactionRun = TTransactionsSchemas['TransactionRun'];

export type TTransactionRun = Model<
  TRawTransactionRun,
  {
    steps: TTestRun[];
    environment?: TEnvironment;
    metadata?: {[key: string]: string};
  }
>;

export type TTransactionRunState = TRawTransactionRun['state'];
