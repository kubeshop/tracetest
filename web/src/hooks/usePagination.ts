import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {useCallback, useState} from 'react';
import {useSearchParams} from 'react-router-dom';

type TParams<P> = P & {
  take?: number;
};

export interface PaginationResponse<T> {
  total: number;
  items: T[];
}

export interface IPagination<T> {
  hasNext: boolean;
  hasPrev: boolean;
  isEmpty: boolean;
  isFetching: boolean;
  isLoading: boolean;
  list: T[];
  loadNext: () => void;
  loadPrev: () => void;
  loadPage: (pg: number) => void;
  search: (query: string) => void;
  page: number;
  total: number;
  take: number;
}

function totalLoaded(list: Array<any>, params: {page: number; query: string}, take: number) {
  const length = list?.length || 0;
  return params.page === 0 ? length : params.page * take + length;
}

const usePagination = <T, P>(
  useGetDataListQuery: UseQuery<any>,
  {take = 20, ...queryParams}: TParams<P>
): IPagination<T> => {
  const [searchParams] = useSearchParams();
  const defaultPage = searchParams.get('page') ? Number(searchParams.get('page')) - 1 : 0;
  const [params, setParams] = useState<{page: number; query: string}>({page: defaultPage, query: ''});

  const {data, isFetching, isLoading} = useGetDataListQuery({
    skip: params.page * take,
    take,
    ...queryParams,
    ...(params.query ? {query: params.query} : {}),
  });

  const list = (data as any)?.items as T[];
  const total = (data as any)?.total as number;

  const loadNext = useCallback(() => {
    setParams(prevParams => ({...prevParams, page: prevParams.page + 1}));
  }, []);

  const loadPrev = useCallback(() => {
    setParams(prevParams => ({...prevParams, page: prevParams.page - 1}));
  }, []);

  const loadPage = useCallback((pg: number) => {
    setParams(prevParams => ({...prevParams, page: pg}));
  }, []);

  const search = (query: string) => {
    setParams({page: 0, query});
  };

  return {
    loadPage,
    hasNext: totalLoaded(list, params, take) < total,
    hasPrev: params.page > 0,
    isEmpty: params.page === 0 && !list?.length && !isLoading && !isFetching,
    isFetching,
    isLoading,
    list: list ?? [],
    loadNext,
    loadPrev,
    search,
    page: params.page + 1,
    total,
    take,
  };
};

export default usePagination;
