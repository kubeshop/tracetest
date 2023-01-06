import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {TDraftTransaction, TTransaction} from 'types/Transaction.types';
import {
  useDeleteTransactionByIdMutation,
  useEditTransactionMutation,
  useLazyGetTransactionVariablesQuery,
  useRunTransactionMutation,
} from 'redux/apis/TraceTest.api';
import {TEnvironmentValue} from 'types/Environment.types';
import {useMissingVariablesModal} from '../../MissingVariablesModal/MissingVariablesModal.provider';

const useTransactionCrud = () => {
  const navigate = useNavigate();
  const [editTransaction, {isLoading: isTransactionEditLoading}] = useEditTransactionMutation();
  const [runTransactionAction, {isLoading: isLoadingRunTransaction}] = useRunTransactionMutation();
  const [deleteTransactionAction] = useDeleteTransactionByIdMutation();
  const isEditLoading = isTransactionEditLoading || isLoadingRunTransaction;
  const [getTransactionVariables, {isLoading: isGetVariablesLoading}] = useLazyGetTransactionVariablesQuery();
  const {selectedEnvironment} = useEnvironment();
  const {onOpen} = useMissingVariablesModal();

  const runTransaction = useCallback(
    async (transaction: TTransaction, runId?: string, environmentId = selectedEnvironment?.id) => {
      const transactionVariables = await getTransactionVariables({
        transactionId: transaction.id,
        version: transaction.version,
        environmentId,
        runId,
      }).unwrap();

      const run = async (variables: TEnvironmentValue[] = []) => {
        const {id} = await runTransactionAction({transactionId: transaction.id, environmentId, variables}).unwrap();

        navigate(`/transaction/${transaction.id}/run/${id}`);
      };

      if (!transactionVariables.hasMissingVariables) run();
      else
        onOpen({
          name: transaction.name,
          testsVariables: transactionVariables.variables,
          onSubmit(variables) {
            run(variables);
          },
        });
    },
    [getTransactionVariables, navigate, onOpen, runTransactionAction, selectedEnvironment?.id]
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
    isGetVariablesLoading,
  };
};

export default useTransactionCrud;
