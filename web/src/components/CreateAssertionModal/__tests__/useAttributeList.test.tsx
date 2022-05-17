import {renderHook} from '@testing-library/react-hooks';
import {TestingModels} from '../../../utils/TestingModels';
import useAttributeList from '../useAttributeList';

test('useAttributeList', () => {
  const {result} = renderHook(() => useAttributeList(TestingModels.span, []));
  expect(result.current[0].options).toEqual([]);
});
