import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {TDraftTransaction, TTransaction} from 'types/Transaction.types';
import {
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useRunTransactionMutation,
} from 'redux/apis/TraceTest.api';
import {TEnvironmentValue} from 'types/Environment.types';
import {useMissingVariablesModal} from 'providers/MissingVariablesModal/MissingVariablesModal.provider';
import {RunErrorTypes, TRunError} from 'types/TestRun.types';

const useTransactionCrud = () => {
  const navigate = useNavigate();
  const [editTransaction, {isLoading: isTransactionEditLoading}] = useEditTransactionMutation();
  const [runTransactionAction, {isLoading: isLoadingRunTransaction}] = useRunTransactionMutation();
  const [deleteTransactionAction] = useDeleteTransactionByIdMutation();
  const isEditLoading = isTransactionEditLoading || isLoadingRunTransaction;
  const {selectedEnvironment} = useEnvironment();
  const {onOpen} = useMissingVariablesModal();

  const runTransaction = useCallback(
    async (transaction: TTransaction, runId?: string, environmentId = selectedEnvironment?.id) => {
      const run = async (variables: TEnvironmentValue[] = []) => {
        try {
          const {id} = await runTransactionAction({transactionId: transaction.id, environmentId, variables}).unwrap();

          navigate(`/transaction/${transaction.id}/run/${id}`);
        } catch (error) {
          const {type, missingVariables} = error as TRunError;
          if (type === RunErrorTypes.MissingVariables)
            onOpen({
              name: transaction.name,
              missingVariables,
              onSubmit(missing) {
                run(missing);
              },
            });
          else throw error;
        }
      };

      run();
    },
    [navigate, onOpen, runTransactionAction, selectedEnvironment?.id]
  );

  const edit = useCallback(
    async (transaction: TTransaction, draft: TDraftTransaction) => {
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
