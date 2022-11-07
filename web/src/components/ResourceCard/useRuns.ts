import {UseLazyQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {useCallback, useState} from 'react';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';

type TParams<P> = P & {
  take: number;
};

export interface IRuns<T> {
  isCollapsed: boolean;
  isLoading: boolean;
  list: T[];
  onClick(): void;
}

const useRuns = <T, P>(useLazyQuery: UseLazyQuery<any>, queryParams: TParams<P>): IRuns<T> => {
  const [isCollapsed, setIsCollapsed] = useState(true);
  const [query, {data, isLoading}] = useLazyQuery();
  const list = ((data as any)?.items as T[]) ?? [];

  const onClick = useCallback(() => {
    if (!isCollapsed) {
      setIsCollapsed(true);
      return;
    }

    setIsCollapsed(false);
    TestAnalyticsService.onTestCardCollapse();
    if (list.length > 0) {
      return;
    }
    query(queryParams);
  }, [query, isCollapsed, list.length, queryParams]);

  return {
    isCollapsed,
    isLoading,
    list,
    onClick,
  };
};

export default useRuns;
