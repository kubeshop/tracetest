import {render} from 'test-utils';
import SpanMock from '../../../models/__mocks__/Span.mock';
import SpanDetail from '../SpanDetail';

test('Layout', () => {
  const {getByText} = render(<SpanDetail span={SpanMock.model()} />);

  expect(getByText('All')).toBeTruthy();
});
