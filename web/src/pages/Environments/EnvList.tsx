import InfiniteScroll from 'components/InfiniteScroll';
import useInfiniteScroll from 'hooks/useInfiniteScroll';
import {IEnvironment, useGetEnvListQuery} from '../../redux/apis/TraceTest.api';
import {EnvironmentCard} from './EnvironmentCard';
import * as S from './Envs.styled';
import NoResults from './NoResults';

interface IProps {
  query: string;
  openDialog: (mode: boolean) => void;
  setEnvironment: (mode: IEnvironment) => void;
}

const EnvList = ({query, openDialog, setEnvironment}: IProps) => {
  const {list, isLoading, loadMore, hasMore} = useInfiniteScroll<IEnvironment, {query: string}>(useGetEnvListQuery, {
    query,
  });
  return (
    <InfiniteScroll
      loadMore={loadMore}
      isLoading={isLoading}
      hasMore={hasMore}
      shouldTrigger={Boolean(list.length)}
      emptyComponent={<NoResults />}
    >
      <S.TestListContainer data-cy="test-list">
        {list?.map(environment => (
          <EnvironmentCard
            key={environment.name}
            environment={environment}
            openDialog={openDialog}
            setEnvironment={setEnvironment}
          />
        ))}
      </S.TestListContainer>
    </InfiniteScroll>
  );
};

export default EnvList;
