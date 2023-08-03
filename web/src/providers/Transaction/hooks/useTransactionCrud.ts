import {useCallback} from 'react';
import {TVariableSetValue} from 'models/VariableSet.model';
import RunError from 'models/RunError.model';
import Transaction from 'models/Transaction.model';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useVariableSet} from 'providers/VariableSet/VariableSet.provider';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import {
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useLazyGetTransactionVersionByIdQuery,
  useRunTransactionMutation,
} from 'redux/apis/Tracetest';
import {RunErrorTypes} from 'types/TestRun.types';
import {TDraftTransaction} from 'types/Transaction.types';

const useTransactionCrud = () => {
  const {navigate} = useDashboard();
  const [editTransaction, {isLoading: isTransactionEditLoading}] = useEditTransactionMutation();
  const [runTransactionAction, {isLoading: isLoadingRunTransaction}] = useRunTransactionMutation();
  const [getTransaction] = useLazyGetTransactionVersionByIdQuery();
  const [deleteTransactionAction] = useDeleteTransactionByIdMutation();
  const isEditLoading = isTransactionEditLoading || isLoadingRunTransaction;
  const {selectedVariableSet} = useVariableSet();
  const {onOpen} = useMissingVariablesModal();

  const runTransaction = useCallback(
    async (transaction: Transaction, runId?: string, variableSetId = selectedVariableSet?.id) => {
      const {fullSteps: testList} = await getTransaction({
        transactionId: transaction.id,
        version: transaction.version,
      }).unwrap();

      const run = async (variables: TVariableSetValue[] = []) => {
        try {
          const {id} = await runTransactionAction({transactionId: transaction.id, variableSetId, variables}).unwrap();

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
    [getTransaction, navigate, onOpen, runTransactionAction, selectedVariableSet?.id]
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
