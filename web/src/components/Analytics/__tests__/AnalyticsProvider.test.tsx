import {render} from 'test-utils';
import AnalyticsProvider from '../AnalyticsProvider';

describe('AnalyticsProvider', () => {
  it('should render the Provider with its children', () => {
    const {getByTestId} = render(
      <AnalyticsProvider>
        <h2 data-cy="sample">Children</h2>
      </AnalyticsProvider>
    );

    expect(getByTestId('sample')).toBeInTheDocument();
  });
});
