import {createContext, useContext, useMemo} from 'react';
import {useGetTransactionRunByIdQuery} from 'redux/apis/TraceTest.api';
import {TTransactionRun} from 'types/TransactionRun.types';
import TransactionProvider from '../Transaction/Transaction.provider';

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

  return transactionRun ? (
    <TransactionProvider transactionId={transactionId} version={transactionRun.version}>
      <Context.Provider value={value}>{children}</Context.Provider>
    </TransactionProvider>
  ) : null;
};

export default TransactionRunProvider;
