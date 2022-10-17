import {render} from 'test-utils';
import TransactionHeader from '../TransactionHeader';

test('TransactionHeader', () => {
  const {getByTestId} = render(<TransactionHeader onBack={jest.fn()} />);
  expect(getByTestId('transaction-details-name')).toBeInTheDocument();
});
