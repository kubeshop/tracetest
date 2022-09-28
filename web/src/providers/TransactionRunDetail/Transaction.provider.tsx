import {createContext, useContext, useMemo} from 'react';
import {useGetTransactionByIdQuery} from 'redux/apis/TraceTest.api';
import {ITransaction} from './ITransaction';

interface IContext {
  transaction?: ITransaction;
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

const TransactionProvider = ({children, transactionId, runId, version = 0}: IProps) => {
  const {data: transaction} = useGetTransactionByIdQuery({transactionId, runId});
  const value = useMemo<IContext>(() => ({transaction, version}), [transaction, version]);

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default TransactionProvider;
