import {createContext, useContext, useMemo} from 'react';
import {useGetTransactionRunByIdQuery} from '../../redux/apis/TraceTest.api';
import {TTransaction} from '../../types/Transaction.types';

interface IContext {
  transaction?: TTransaction;
  version: number;
}

export const Context = createContext<IContext>({
  transaction: undefined,
  version: 0,
});

interface IProps {
  transactionId: string;
  runId: string;
  version?: number;
  children: React.ReactNode;
}

export const useTransaction = () => useContext(Context);

const TransactionRunDetailProvider = ({children, transactionId, runId, version = 0}: IProps) => {
  const {data: transaction} = useGetTransactionRunByIdQuery({transactionId, runId});
  const value = useMemo<IContext>(() => ({transaction, version}), [transaction, version]);

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TransactionRunDetailProvider;
