import {useCallback} from 'react';
import {TTest} from '../../types/Test.types';
import {useDeleteTestByIdMutation} from '../../redux/apis/TraceTest.api';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';

export function useMenuDeleteCallback(): (test: TTest) => void {
  const [deleteTestMutation] = useDeleteTestByIdMutation();

  return useCallback(
    ({id: testId}: TTest) => {
      TestAnalyticsService.onDeleteTest();
      deleteTestMutation({testId});
    },
    [deleteTestMutation]
  );
}
