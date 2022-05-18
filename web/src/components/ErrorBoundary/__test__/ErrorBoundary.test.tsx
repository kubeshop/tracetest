import {render} from '@testing-library/react';
import ErrorBoundary from '../ErrorBoundary';

test('ErrorBoundary', () => {
  const result = render(<ErrorBoundary error={new Error('cannot find value of undefined')} />);
  expect(result.container).toMatchSnapshot();
});
