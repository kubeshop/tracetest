import {render} from '@testing-library/react';
import TestCard from '../TestCard';
import {TestingModels} from '../../../utils/TestingModels';

test('SpanAttributesTable', () => {
  const onDelete = jest.fn();
  const onRunTest = jest.fn();
  const onClick = jest.fn();
  const result = render(
    <TestCard onDelete={onDelete} onRunTest={onRunTest} test={TestingModels.test} onClick={onClick} />
  );
  expect(result.container).toMatchSnapshot();
});
