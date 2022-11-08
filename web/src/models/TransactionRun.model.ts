import {TRawTransactionRun, TTransactionRun} from 'types/Transaction.types';
import TransactionRunResult from './TransactionRunResult';
import Environment from './Environment.model';

const TransactionRunModel = ({id = '', environment = {}, results = []}: TRawTransactionRun): TTransactionRun => {
  return {
    id,
    environment: Environment(environment),
    results: results.map(result => TransactionRunResult(result)),
  };
};

export default TransactionRunModel;
