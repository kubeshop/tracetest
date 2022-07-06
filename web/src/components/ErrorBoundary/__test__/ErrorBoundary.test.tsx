import {render} from 'test-utils';
import ErrorBoundary from '../ErrorBoundary';

test('ErrorBoundary', () => {
  const errorMsg = 'cannot find value of undefined';
  const {getByText} = render(<ErrorBoundary error={new Error(errorMsg)} />);

  expect(getByText('Something went wrong!')).toBeTruthy();
});
