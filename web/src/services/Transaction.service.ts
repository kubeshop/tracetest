import {TDraftTransaction} from 'types/Transaction.types';
import Transaction from 'models/Transaction.model';

const TransactionService = () => ({
  getRawFromDraft(draftTransaction: TDraftTransaction) {
    return {...draftTransaction, steps: draftTransaction.steps?.map(step => ({id: step}))};
  },

  getInitialValues(transaction: Transaction): TDraftTransaction {
    return {
      ...transaction,
      steps: transaction.steps.map(step => step.id),
    };
  },
});

export default TransactionService();
