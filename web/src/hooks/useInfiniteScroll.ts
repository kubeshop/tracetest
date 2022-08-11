import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {isArray} from 'lodash';
import {useCallback, useEffect, useMemo, useState} from 'react';

type TUseInfiniteScrollParams<P> = P & {
  take?: number;
};

export interface InfiniteScrollModel<T> {
  isLoading: boolean;
  localPage: number;
  loadMore: () => void;
  hasMore: boolean;
  isFetching: boolean;
  refresh: () => void;
  list: T[];
}

const useInfiniteScroll = <T, P>(
  useGetDataListQuery: UseQuery<any>,
  {take = 20, ...queryParams}: TUseInfiniteScrollParams<P>
): InfiniteScrollModel<T> => {
  const [localPage, setLocalPage] = useState(0);
  const [list, setList] = useState<T[]>([]);
  const [lastCount, setLastCount] = useState(0);

  const {data, isLoading, isFetching} = useGetDataListQuery({
    skip: localPage * take,
    take,
    ...queryParams,
  });

  const currentList = data as T[];

  const hasMore = useMemo(() => lastCount === take, [lastCount, take]);

  useEffect(() => {
    if (isArray(currentList)) {
      if (localPage === 0) {
        setList(currentList);
      } else if (localPage > 0) {
        setList(prevList => [...prevList, ...currentList]);
      }

      setLastCount(currentList.length);
    }
  }, [currentList, localPage]);

  const refresh = useCallback(() => {
    setLocalPage(1);
  }, []);

  const loadMore = useCallback(() => {
    setLocalPage(page => page + 1);
  }, []);

  return {
    list,
    localPage,
    isLoading,
    isFetching,
    hasMore,
    loadMore,
    refresh,
  };
};

export default useInfiniteScroll;
