import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {isArray} from 'lodash';
import {useCallback, useEffect, useMemo, useState} from 'react';

type TUseInfiniteScrollParams<P> = P & {
  take?: number;
};

const useInfiniteScroll = <T, P>(
  useGetDataListQuery: UseQuery<any>,
  {take = 20, ...queryParams}: TUseInfiniteScrollParams<P>
) => {
  const [localPage, setLocalPage] = useState(0);
  const [list, setList] = useState<T[]>([]);
  const [lastCount, setLastCount] = useState(0);
  console.log('@@', {localPage, take});

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
  }, [currentList]);

  const refresh = useCallback(() => {
    setLocalPage(1);
  }, []);

  const loadMore = useCallback(() => {
    console.log('@@loadMore');
    setLocalPage(page => page + 1);
  }, []);

  console.log('@@', {
    list,
    localPage,
    isLoading,
    isFetching,
    hasMore,
    loadMore,
    refresh,
    lastCount,
  });

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
