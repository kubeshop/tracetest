import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useMemo} from 'react';
import {useNavigate} from 'react-router-dom';

import {useGetTransactionByIdQuery, useDeleteTransactionByIdMutation} from 'redux/apis/TraceTest.api';
import {TTransaction} from 'types/Transaction.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  isError: boolean;
  isLoading: boolean;
  isLoadingRun: boolean;
  onDelete(id: string, name: string): void;
  onRun(): void;
  transaction: TTransaction;
}

export const Context = createContext<IContext>({
  isError: false,
  isLoading: false,
  isLoadingRun: false,
  onDelete: noop,
  onRun: noop,
  transaction: {} as TTransaction,
});

interface IProps {
  children: ReactNode;
  transactionId: string;
}

export const useTransaction = () => useContext(Context);

const TransactionProvider = ({children, transactionId}: IProps) => {
  const {data: transaction, isLoading, isError} = useGetTransactionByIdQuery({transactionId});
  const [deleteTransaction] = useDeleteTransactionByIdMutation();
  const {onOpen} = useConfirmationModal();
  const navigate = useNavigate();

  const onRun = useCallback(() => {
    console.log('onRun');
  }, []);

  const onDelete = useCallback(
    (id: string, name: string) => {
      function onConfirmation() {
        deleteTransaction({transactionId: id});
        navigate('/');
      }

      onOpen(`Are you sure you want to delete “${name}”?`, onConfirmation);
    },
    [navigate, onOpen]
  );

  const value = useMemo<IContext>(
    () => ({
      isError,
      isLoading,
      isLoadingRun: false,
      onDelete,
      onRun,
      transaction: transaction!,
    }),
    [isError, isLoading, onDelete, onRun, transaction]
  );

  return transaction ? (
    <Context.Provider value={value}>{children}</Context.Provider>
  ) : (
    <div data-cy="loading-transaction" />
  );
};

export default TransactionProvider;
