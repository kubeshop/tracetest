import {TDraftTransaction} from 'types/Transaction.types';
import Transaction, {TRawTransaction} from 'models/Transaction.model';

const TransactionService = () => ({
  getRawFromDraft(draftTransaction: TDraftTransaction): TRawTransaction {
    return {
      spec: {...draftTransaction, fullSteps: draftTransaction.steps?.map(step => ({id: step}))},
      type: 'Transaction',
    };
  },

  getInitialValues(transaction: Transaction): TDraftTransaction {
    return {
      ...transaction,
      steps: transaction.fullSteps.map(step => step.id),
    };
  },
});

export default TransactionService();
