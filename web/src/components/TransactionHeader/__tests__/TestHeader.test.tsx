import {render} from 'test-utils';
import TransactionMock from '../../../models/__mocks__/Transaction.mock';
import TransactionHeader from '../TransactionHeader';

test('TransactionHeader', () => {
  const {getByTestId} = render(<TransactionHeader onBack={jest.fn()} transaction={TransactionMock.model()} />);
  expect(getByTestId('transaction-details-name')).toBeInTheDocument();
});
