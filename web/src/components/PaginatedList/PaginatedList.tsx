import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {QueryDefinition} from '@reduxjs/toolkit/query';
import {ReactElement} from 'react';
import {PaginationResponse} from '../../hooks/usePagination';
import Pagination from '../Pagination';
import Empty from './Empty';
import Loading from './Loading';
import * as S from './PaginatedList.styled';

type TParams<P> = P & {
  take?: number;
};

interface IProps<T, P> {
  dataCy?: string;
  params: TParams<P>;
  query: UseQuery<QueryDefinition<P, any, any, PaginationResponse<T>>>;
  itemComponent({item}: {item: T}): ReactElement;
}

interface IId {
  id: string;
}

const PaginatedList = <T extends IId, P>({dataCy, itemComponent: ItemComponent, params, query}: IProps<T, P>) => {
  return (
    <Pagination<T, P>
      emptyComponent={<Empty />}
      loadingComponent={<Loading />}
      query={query}
      defaultParameters={params}
    >
      {pagination => (
        <S.ListContainer data-cy={dataCy}>
          {pagination.list.map(item => (
            <ItemComponent item={item} key={item.id} />
          ))}
        </S.ListContainer>
      )}
    </Pagination>
  );
};

export default PaginatedList;
