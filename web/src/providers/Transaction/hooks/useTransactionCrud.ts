import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {TDraftTransaction} from 'types/Transaction.types';
import {
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useLazyGetTransactionVersionByIdQuery,
  useRunTransactionMutation,
} from 'redux/apis/TraceTest.api';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import {RunErrorTypes} from 'types/TestRun.types';
import Transaction from 'models/Transaction.model';
import {TEnvironmentValue} from 'models/Environment.model';
import RunError from 'models/RunError.model';

const useTransactionCrud = () => {
  const navigate = useNavigate();
  const [editTransaction, {isLoading: isTransactionEditLoading}] = useEditTransactionMutation();
  const [runTransactionAction, {isLoading: isLoadingRunTransaction}] = useRunTransactionMutation();
  const [getTransaction] = useLazyGetTransactionVersionByIdQuery();
  const [deleteTransactionAction] = useDeleteTransactionByIdMutation();
  const isEditLoading = isTransactionEditLoading || isLoadingRunTransaction;
  const {selectedEnvironment} = useEnvironment();
  const {onOpen} = useMissingVariablesModal();

  const runTransaction = useCallback(
    async (transaction: Transaction, runId?: string, environmentId = selectedEnvironment?.id) => {
      const {fullSteps: testList} = await getTransaction({
        transactionId: transaction.id,
        version: transaction.version,
      }).unwrap();

      const run = async (variables: TEnvironmentValue[] = []) => {
        try {
          const {id} = await runTransactionAction({transactionId: transaction.id, environmentId, variables}).unwrap();

          navigate(`/transaction/${transaction.id}/run/${id}`);
        } catch (error) {
          const {type, missingVariables} = error as RunError;
          if (type === RunErrorTypes.MissingVariables)
            onOpen({
              name: transaction.name,
              missingVariables,
              testList,
              onSubmit(missing) {
                run(missing);
              },
            });
          else throw error;
        }
      };

      run();
    },
    [getTransaction, navigate, onOpen, runTransactionAction, selectedEnvironment?.id]
  );

  const edit = useCallback(
    async (transaction: Transaction, draft: TDraftTransaction) => {
      await editTransaction({transactionId: transaction.id, transaction: draft}).unwrap();

      runTransaction(transaction);
    },
    [editTransaction, runTransaction]
  );

  const deleteTransaction = useCallback(
    (transactionId: string) => {
      deleteTransactionAction({transactionId});

      navigate('/');
    },
    [deleteTransactionAction, navigate]
  );

  return {
    edit,
    runTransaction,
    deleteTransaction,
    isEditLoading,
  };
};

export default useTransactionCrud;
