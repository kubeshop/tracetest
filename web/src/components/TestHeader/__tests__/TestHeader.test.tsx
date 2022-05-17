import {render} from '@testing-library/react';
import {TestState} from '../../../constants/TestRunResult.constants';
import {TestingModels} from '../../../utils/TestingModels';
import TestHeader from '../TestHeader';

test('SpanAttributesTable', () => {
  const result = render(<TestHeader onBack={jest.fn()} test={TestingModels.test} testState={TestState.CREATED} />);
  expect(result.container).toMatchSnapshot();
});
