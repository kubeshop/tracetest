import {createContext, useContext, useMemo} from 'react';
import {useGetTransactionRunByIdQuery} from 'redux/apis/TraceTest.api';
import {TTransactionRun} from 'types/Transaction.types';

interface IContext {
  transactionRun?: TTransactionRun;
}

export const Context = createContext<IContext>({
  transactionRun: undefined,
});

interface IProps {
  transactionId: string;
  runId: string;
  children: React.ReactNode;
}

export const useTransactionRun = () => useContext(Context);

const TransactionRunProvider = ({children, transactionId, runId}: IProps) => {
  const {data: transactionRun} = useGetTransactionRunByIdQuery({transactionId, runId});

  const value = useMemo<IContext>(() => ({transactionRun}), [transactionRun]);

  return transactionRun ? <Context.Provider value={value}>{children}</Context.Provider> : null;
};

export default TransactionRunProvider;
