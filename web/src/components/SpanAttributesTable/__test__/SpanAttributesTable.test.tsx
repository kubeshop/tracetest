import {render, waitFor} from '@testing-library/react';
import SpanAttributesTable from '../SpanAttributesTable';

test('SpanAttributesTable', async () => {
  const {container, getByText} = render(<SpanAttributesTable spanAttributesList={[]} />);

  await waitFor(() => getByText('No Data'));
  expect(container).toMatchSnapshot();
});
