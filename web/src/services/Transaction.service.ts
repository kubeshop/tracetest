import {TDraftTransaction, TTransaction} from 'types/Transaction.types';

const TransactionService = () => ({
  getRawFromDraft(draftTransaction: TDraftTransaction) {
    return {...draftTransaction, steps: draftTransaction.steps?.map(step => ({id: step}))};
  },

  getInitialValues(transaction: TTransaction): TDraftTransaction {
    return {
      ...transaction,
      steps: transaction.steps.map(step => step.id),
    };
  },
});

export default TransactionService();
