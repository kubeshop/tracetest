import {render} from '@testing-library/react';
import ErrorBoundary from '../ErrorBoundary';

test('ErrorBoundary', () => {
  const result = render(<ErrorBoundary error={{message: '', stack: '', name: 'sdfk'}} />);
  expect(result.container).toMatchSnapshot();
});
