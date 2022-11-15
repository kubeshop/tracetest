import {render} from 'test-utils';
import TransactionRunModel from '../../../models/TransactionRun.model';
import TransactionMock from '../../../models/__mocks__/Transaction.mock';
import TransactionHeader from '../TransactionHeader';

test('TransactionHeader', () => {
  const {getByTestId} = render(
    <TransactionHeader
      onBack={jest.fn()}
      transaction={TransactionMock.model()}
      transactionRun={TransactionRunModel({})}
    />
  );
  expect(getByTestId('transaction-details-name')).toBeInTheDocument();
});
