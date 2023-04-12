import {Pagination as PG} from 'antd';
import {useNavigate} from 'react-router-dom';
import {ReactNode, useCallback, useEffect} from 'react';
import {IPagination} from 'hooks/usePagination';

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
  const navigate = useNavigate();
  const handleNextPage = useCallback(() => {
    if (isLoading || !hasNext) return;
    loadNext();
  }, [hasNext, isLoading, loadNext]);

  const handlePrevPage = useCallback(() => {
    if (isLoading || !hasPrev) return;
    loadPrev();
  }, [hasPrev, isLoading, loadPrev]);

  useEffect(() => {
    navigate(
      {
        search: `page=${page}`,
      },
      {
        replace: true,
      }
    );
  }, [navigate, page]);

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
            showSizeChanger={false}
          />
        </S.Footer>
      )}
    </>
  );
};

export default Pagination;
