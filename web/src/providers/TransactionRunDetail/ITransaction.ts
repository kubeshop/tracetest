import {TTest} from '../../types/Test.types';

interface TransactionStep extends TTest {
  result: 'success' | 'fail' | 'running';
}

export interface ITransaction {
  id: string;
  name: string;
  description: string;
  version: number;
  steps: TransactionStep[];
  env: Record<string, string>;
}
