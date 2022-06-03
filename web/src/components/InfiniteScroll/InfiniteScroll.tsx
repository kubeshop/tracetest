import {useEffect, ReactNode} from 'react';
import {useInView} from 'react-intersection-observer';

interface IInfiniteScrollProps {
  shouldTrigger: boolean;
  loadMore(): void;
  hasMore: boolean;
  isLoading: boolean;
  emptyComponent?: ReactNode;
}

const InfiniteScroll: React.FC<IInfiniteScrollProps> = ({
  children,
  isLoading,
  hasMore,
  shouldTrigger,
  loadMore,
  emptyComponent,
}) => {
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
      {!isLoading && !shouldTrigger && emptyComponent}
    </>
  );
};

export default InfiniteScroll;
