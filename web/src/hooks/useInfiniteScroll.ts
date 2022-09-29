import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {isArray} from 'lodash';
import {useCallback, useEffect, useMemo, useState} from 'react';

type TUseInfiniteScrollParams<P> = P & {
  take?: number;
};

export interface InfiniteScrollModel<T> {
  hasMore: boolean;
  isLoading: boolean;
  list: T[];
  loadMore: () => void;
  search: (query: string) => void;
}

const useInfiniteScroll = <T, P>(
  useGetDataListQuery: UseQuery<any>,
  {take = 20, ...queryParams}: TUseInfiniteScrollParams<P>
): InfiniteScrollModel<T> => {
  const [list, setList] = useState<T[]>([]);
  const [lastCount, setLastCount] = useState(0);
  const [params, setParams] = useState<{page: number; query: string}>({page: 0, query: ''});

  const hasMore = useMemo(() => lastCount === take, [lastCount, take]);

  const {data, isFetching} = useGetDataListQuery({
    skip: params.page * take,
    take,
    ...queryParams,
    ...(params.query ? {query: params.query} : {}),
  });

  useEffect(() => {
    const currentList = data as T[];
    if (!isArray(currentList)) return;

    setList(prevList => [...prevList, ...currentList]);
    setLastCount(currentList.length);
  }, [data]);

  const loadMore = useCallback(() => {
    setParams(prevParams => ({...prevParams, page: prevParams.page + 1}));
  }, []);

  const search = (query: string) => {
    setList([]);
    setLastCount(0);
    setParams({page: 0, query});
  };

  return {
    hasMore,
    isLoading: isFetching,
    list,
    loadMore,
    search,
  };
};

export default useInfiniteScroll;
