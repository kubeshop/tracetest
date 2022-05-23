import {useCallback} from 'react';
import {TTest} from '../../types/Test.types';
import {useDeleteTestByIdMutation} from '../../redux/apis/TraceTest.api';

export function useMenuDeleteCallback(): (test: TTest) => void {
  const [deleteTestMutation] = useDeleteTestByIdMutation();

  return useCallback(({id: testId}: TTest) => deleteTestMutation({testId}), [deleteTestMutation]);
}
