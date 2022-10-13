import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {ReactElement} from 'react';
import usePagination from '../../hooks/usePagination';
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
  query: UseQuery<any>;
  itemComponent({item}: {item: T}): ReactElement;
}

interface IId {
  id: string;
}

const PaginatedList = <T extends IId, P>({dataCy, itemComponent: ItemComponent, params, query}: IProps<T, P>) => {
  const pagination = usePagination<T, typeof params>(query, params);

  return (
    <Pagination emptyComponent={<Empty />} loadingComponent={<Loading />} {...pagination}>
      <S.ListContainer data-cy={dataCy}>
        {pagination.list.map(item => (
          <ItemComponent item={item} key={item.id} />
        ))}
      </S.ListContainer>
    </Pagination>
  );
};

export default PaginatedList;
