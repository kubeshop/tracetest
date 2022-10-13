import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {QueryDefinition} from '@reduxjs/toolkit/query';
import {useCallback, useState} from 'react';

export type TParams<P> = P & {
  take?: number;
};

export interface PaginationResponse<T> {
  total: number;
  items: T[];
}

export interface IPagination<T> {
  isEmpty: boolean;
  isFetching: boolean;
  isLoading: boolean;
  list: T[];
  loadPage: (pg: number) => void;
  search: (query: string) => void;
  page: number;
  total: number;
  take: number;
}

const usePagination = <T, P>(
  useGetDataListQuery: UseQuery<QueryDefinition<P, any, any, PaginationResponse<T>>>,
  {...queryParams}: TParams<P>
): IPagination<T> => {
  const take = queryParams?.take || 5;
  const [params, setParams] = useState<{page: number; query: string}>({page: 0, query: ''});
  const {data, isFetching, isLoading} = useGetDataListQuery({
    skip: params.page * take,
    take,
    ...queryParams,
    ...(params.query ? {query: params.query} : {}),
  });
  const list = data?.items || [];
  return {
    isFetching,
    isLoading,
    list: list ?? [],
    page: params.page + 1,
    total: data?.total || 0,
    take,
    isEmpty: params.page === 0 && !list?.length && !isLoading && !isFetching,
    search: useCallback(query => setParams({page: 0, query}), [setParams]),
    loadPage: useCallback(pg => {
      setParams(prevParams => ({...prevParams, page: pg}));
    }, []),
  };
};

export default usePagination;
