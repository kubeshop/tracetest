import {render} from '@testing-library/react';
import SpanMock from '../../../models/__mocks__/Span.mock';
import Http from '../components/Http/Http';

test('Http', () => {
  const {getAllByTestId} = render(
    <Http
      onCreateAssertion={jest.fn()}
      span={SpanMock.model()}
    />
  );
  expect(getAllByTestId('span-details-attributes').length).toBe(1);
});
