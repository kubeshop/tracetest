import {render} from '@testing-library/react';
import Http from '../components/Http/Http';
import {TestingModels} from '../../../utils/TestingModels';

test('Http', () => {
  const {getAllByTestId} = render(
    <Http
      onCreateAssertion={jest.fn()}
      span={TestingModels.span}
    />
  );
  expect(getAllByTestId('span-details-attributes').length).toBe(1);
});
