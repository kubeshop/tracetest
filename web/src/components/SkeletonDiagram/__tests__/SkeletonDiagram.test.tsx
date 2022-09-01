import {render} from 'test-utils';
import SkeletonDiagram from '../SkeletonDiagram';

describe('SkeletonDiagram', () => {
  it('should render correctly', () => {
    const {getByTestId} = render(<SkeletonDiagram />);

    expect(getByTestId('skeleton-diagram')).toBeInTheDocument();
  });
});
