import {useEffect, ReactNode} from 'react';
import {useInView} from 'react-intersection-observer';

interface IInfiniteScrollProps {
  children: ReactNode;
  shouldTrigger: boolean;
  loadMore(): void;
  hasMore: boolean;
  isLoading: boolean;
  emptyComponent?: ReactNode;
  loadingComponent?: ReactNode;
}

const InfiniteScroll = ({
  children,
  isLoading,
  hasMore,
  shouldTrigger,
  loadMore,
  emptyComponent,
  loadingComponent,
}: IInfiniteScrollProps) => {
  const {ref, inView} = useInView();

  useEffect(() => {
    if (inView && !isLoading && hasMore) {
      loadMore();
    }
  }, [hasMore, inView, isLoading, loadMore]);

  return (
    <>
      {children}
      {shouldTrigger && <div style={{height: '10px'}} ref={ref} />}
      {isLoading && loadingComponent}
      {!isLoading && !shouldTrigger && emptyComponent}
    </>
  );
};

export default InfiniteScroll;
