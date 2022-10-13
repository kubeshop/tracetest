import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {QueryDefinition} from '@reduxjs/toolkit/query';
import {Pagination as PaginationAntd} from 'antd';
import {Dispatch, ReactNode, SetStateAction, useState} from 'react';
import usePagination, {IPagination, PaginationResponse, TParams} from '../../hooks/usePagination';
import Loading from '../../pages/Home/Loading';
import NoResults from '../../pages/Home/NoResults';

import * as S from './Pagination.styled';

interface IProps<T, P> {
  children: (pagination: IPagination<T>, params: [P, Dispatch<SetStateAction<P>>]) => ReactNode;
  defaultParameters: TParams<P>;
  query: UseQuery<QueryDefinition<P, any, any, PaginationResponse<T>>>;
  emptyComponent?: ReactNode;
  loadingComponent?: ReactNode;
}

const Pagination = <T, P>({
  children,
  defaultParameters,
  query,
  emptyComponent,
  loadingComponent,
}: IProps<T, P>): React.ReactElement => {
  const [parameters, setParameters] = useState<TParams<P>>(defaultParameters);
  const pagination = usePagination<T, P>(query, parameters);
  return (
    <>
      {children(pagination, [parameters, setParameters])}
      {pagination.isLoading && (emptyComponent || <Loading />)}
      {pagination.isEmpty && (loadingComponent || <NoResults />)}
      {!pagination.isEmpty && (
        <S.Footer>
          <PaginationAntd
            showSizeChanger={false}
            onChange={pg => pagination.loadPage(pg - 1)}
            pageSize={pagination.take}
            current={pagination.page}
            total={pagination.total}
          />
        </S.Footer>
      )}
    </>
  );
};

export default Pagination;
