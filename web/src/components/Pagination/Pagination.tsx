import {LeftOutlined, RightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {ReactNode} from 'react';

import * as S from './Pagination.styled';

interface IProps {
  children: ReactNode;
  emptyComponent?: ReactNode;
  hasNext: boolean;
  hasPrev: boolean;
  isEmpty: boolean;
  isFetching: boolean;
  isLoading: boolean;
  loadingComponent?: ReactNode;
  loadNext(): void;
  loadPrev(): void;
}

const Pagination = ({
  children,
  emptyComponent,
  hasNext,
  hasPrev,
  isEmpty,
  isFetching,
  isLoading,
  loadingComponent,
  loadNext,
  loadPrev,
}: IProps) => {
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
          <div>
            {hasPrev && (
              <Button ghost loading={isFetching} onClick={handlePrevPage} type="primary">
                <LeftOutlined /> Prev page
              </Button>
            )}
          </div>
          <div>
            {hasNext && (
              <Button disabled={!hasNext} ghost loading={isFetching} onClick={handleNextPage} type="primary">
                Next page <RightOutlined />
              </Button>
            )}
          </div>
        </S.Footer>
      )}
    </>
  );
};

export default Pagination;
