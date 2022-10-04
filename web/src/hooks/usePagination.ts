import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {useCallback, useState} from 'react';

type TParams<P> = P & {
  take?: number;
};

interface IPagination<T> {
  hasNext: boolean;
  hasPrev: boolean;
  isEmpty: boolean;
  isFetching: boolean;
  isLoading: boolean;
  list: T[];
  loadNext: () => void;
  loadPrev: () => void;
  search: (query: string) => void;
}

const usePagination = <T, P>(
  useGetDataListQuery: UseQuery<any>,
  {take = 20, ...queryParams}: TParams<P>
): IPagination<T> => {
  const [params, setParams] = useState<{page: number; query: string}>({page: 0, query: ''});

  const {data, isFetching, isLoading} = useGetDataListQuery({
    skip: params.page * take,
    take,
    ...queryParams,
    ...(params.query ? {query: params.query} : {}),
  });

  const list = data as T[];

  const loadNext = useCallback(() => {
    setParams(prevParams => ({...prevParams, page: prevParams.page + 1}));
  }, []);

  const loadPrev = useCallback(() => {
    setParams(prevParams => ({...prevParams, page: prevParams.page - 1}));
  }, []);

  const search = (query: string) => {
    setParams({page: 0, query});
  };

  return {
    hasNext: list?.length === take,
    hasPrev: params.page > 0,
    isEmpty: params.page === 0 && !list?.length && !isLoading && !isFetching,
    isFetching,
    isLoading,
    list: list ?? [],
    loadNext,
    loadPrev,
    search,
  };
};

export default usePagination;
