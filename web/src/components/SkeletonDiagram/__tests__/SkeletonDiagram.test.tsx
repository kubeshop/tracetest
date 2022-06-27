import noop from 'lodash/noop';
import {render} from 'test-utils';
import SkeletonDiagram from '../SkeletonDiagram';

describe('SkeletonDiagram', () => {
  it('should render correctly', () => {
    const {getByTestId} = render(<SkeletonDiagram onClearAffectedSpans={noop} onClearSelectedSpan={noop} />);

    expect(getByTestId('skeleton-diagram')).toBeInTheDocument();
  });
});
