import {render} from '@testing-library/react';
import SkeletonTable from '../SkeletonTable';

test('SkeletonTable', () => {
  const result = render(
    <SkeletonTable loading={false}>
      <h2>Whatever</h2>
    </SkeletonTable>
  );
  expect(result.container).toMatchSnapshot();
});
