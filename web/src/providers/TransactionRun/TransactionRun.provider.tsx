import {createContext, useContext, useMemo} from 'react';
import {useGetTransactionRunByIdQuery} from 'redux/apis/Tracetest';
import TransactionRun from 'models/TransactionRun.model';
import TransactionProvider from '../Transaction/Transaction.provider';

interface IContext {
  transactionRun: TransactionRun;
}

export const Context = createContext<IContext>({
  transactionRun: {} as TransactionRun,
});

interface IProps {
  transactionId: string;
  runId: string;
  children: React.ReactNode;
}

export const useTransactionRun = () => useContext(Context);

const TransactionRunProvider = ({children, transactionId, runId}: IProps) => {
  const {data: transactionRun} = useGetTransactionRunByIdQuery({transactionId, runId});
  const value = useMemo<IContext>(() => ({transactionRun: transactionRun!}), [transactionRun]);

  return transactionRun ? (
    <TransactionProvider transactionId={transactionId} version={transactionRun.version}>
      <Context.Provider value={value}>{children}</Context.Provider>
    </TransactionProvider>
  ) : null;
};

export default TransactionRunProvider;
