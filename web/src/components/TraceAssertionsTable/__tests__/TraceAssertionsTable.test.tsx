import {act, render} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import {store} from '../../../redux/store';
import TraceAssertionsTable from '../TraceAssertionsTable';

test('Layout', () => {
  act(() => {
    const result = render(
      <MemoryRouter>
        <Provider store={store}>
          <TraceAssertionsTable
            assertionResult={{
              assertion: {assertionId: '', selectors: undefined, spanAssertions: undefined},
              spanListAssertionResult: [],
            }}
            onSpanSelected={jest.fn()}
          >
            <h2>This</h2>
          </TraceAssertionsTable>
        </Provider>
      </MemoryRouter>
    );
  });
  // const input = screen.getByText('This');
  expect(result.container).toMatchSnapshot();
  // expect(input).toBeTruthy();
});
