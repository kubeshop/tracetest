import {render} from '@testing-library/react';
import SpanAttributesTable from '../SpanAttributesTable';

test('SpanAttributesTable', () => {
  const result = render(<SpanAttributesTable spanAttributesList={[]} />);
  expect(result.container).toMatchSnapshot();
});
