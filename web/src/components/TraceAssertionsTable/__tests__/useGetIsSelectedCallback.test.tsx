import {renderHook} from '@testing-library/react-hooks';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import {TestingModels} from '../../../utils/TestingModels';
import {useGetIsSelectedCallback} from '../useGetIsSelectedCallback';

test('useGetIsSelectedCallback', () => {
  const {result} = renderHook(() => useGetIsSelectedCallback(), {wrapper: ReduxWrapperProvider});
  expect(result.current(TestingModels.spanId)).toStrictEqual(false);
});

test('useGetIsSelectedCallback', () => {
  const {result} = renderHook(() => useGetIsSelectedCallback(), {wrapper: ReduxWrapperProvider});
  expect(result.current(TestingModels.spanId)).toStrictEqual(false);
});
