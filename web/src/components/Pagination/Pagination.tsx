import {Pagination as PG} from 'antd';
import {ReactNode} from 'react';
import {IPagination} from '../../hooks/usePagination';

import * as S from './Pagination.styled';

interface IProps<T> extends IPagination<T> {
  children: ReactNode;
  emptyComponent?: ReactNode;
  loadingComponent?: ReactNode;
}

const Pagination = <T extends any>({
  children,
  emptyComponent,
  hasNext,
  hasPrev,
  isEmpty,
  isLoading,
  loadingComponent,
  loadNext,
  loadPrev,
  total,
  page,
  take,
  loadPage,
}: IProps<T>) => {
  const handleNextPage = () => {
    if (isLoading || !hasNext) return;
    loadNext();
  };

  const handlePrevPage = () => {
    if (isLoading || !hasPrev) return;
    loadPrev();
  };

  return (
    <>
      {children}
      {isLoading && loadingComponent}
      {isEmpty && emptyComponent}
      {!isEmpty && (
        <S.Footer>
          <PG
            onChange={pg => {
              if (page - 1 === pg) {
                handlePrevPage();
                return;
              }
              if (page + 1 === pg) {
                handleNextPage();
                return;
              }
              loadPage(pg - 1);
            }}
            pageSize={take}
            current={page}
            total={total}
          />
        </S.Footer>
      )}
    </>
  );
};

export default Pagination;
