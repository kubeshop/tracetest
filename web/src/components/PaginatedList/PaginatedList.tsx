import {UseQuery} from '@reduxjs/toolkit/dist/query/react/buildHooks';
import {ReactElement} from 'react';

import usePagination from 'hooks/usePagination';
import Pagination from 'components/Pagination';
import Empty from './Empty';
import Loading from './Loading';
import * as S from './PaginatedList.styled';

type TParams<P> = P & {
  take?: number;
};

interface IProps<T, P> {
  itemComponent({item}: {item: T}): ReactElement;
  params: TParams<P>;
  query: UseQuery<any>;
}

interface IId {
  id: string;
}

const PaginatedList = <T extends IId, P>({itemComponent: ItemComponent, params, query}: IProps<T, P>) => {
  const {hasNext, hasPrev, isEmpty, isFetching, isLoading, list, loadNext, loadPrev} = usePagination<T, typeof params>(
    query,
    params
  );

  return (
    <Pagination
      emptyComponent={<Empty />}
      hasNext={hasNext}
      hasPrev={hasPrev}
      isEmpty={isEmpty}
      isFetching={isFetching}
      isLoading={isLoading}
      loadingComponent={<Loading />}
      loadNext={loadNext}
      loadPrev={loadPrev}
    >
      <S.ListContainer>
        {list.map(item => (
          <ItemComponent item={item} key={item.id} />
        ))}
      </S.ListContainer>
    </Pagination>
  );
};

export default PaginatedList;
