import {useCallback} from 'react';
import {ITest} from '../../types/Test.types';
import {useDeleteTestByIdMutation} from '../../redux/apis/Test.api';

export function useMenuDeleteCallback(): (test: ITest) => void {
  const [deleteTestMutation] = useDeleteTestByIdMutation();

  return useCallback(({testId}: ITest) => deleteTestMutation(testId), [deleteTestMutation]);
}
