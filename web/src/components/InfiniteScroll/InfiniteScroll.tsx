import {useEffect} from 'react';
import {useInView} from 'react-intersection-observer';

interface IInfiniteScrollProps {
  shouldTrigger: boolean;
  loadMore(): void;
  hasMore: boolean;
  isLoading: boolean;
}

const InfiniteScroll: React.FC<IInfiniteScrollProps> = ({children, isLoading, hasMore, shouldTrigger, loadMore}) => {
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
    </>
  );
};

export default InfiniteScroll;
