import {Dispatch, SetStateAction} from 'react';
import Pagination from '../../components/Pagination';
import usePagination from '../../hooks/usePagination';
import {useGetEnvListQuery} from '../../redux/apis/TraceTest.api';
import Loading from '../Home/Loading';
import NoResults from '../Home/NoResults';
import * as S from './Environment.styled';
import {EnvironmentCard} from './EnvironmentCard';
import {IEnvironment} from './IEnvironment';

interface IProps {
  query: string;
  setIsFormOpen: Dispatch<SetStateAction<boolean>>;
  setEnvironment: Dispatch<SetStateAction<IEnvironment | undefined>>;
}

const EnvironmentList = ({query, setEnvironment, setIsFormOpen}: IProps) => {
  const {hasNext, hasPrev, isEmpty, isFetching, isLoading, list, loadNext, loadPrev} = usePagination<
    IEnvironment,
    {query: string}
  >(useGetEnvListQuery, {query});
  return (
    <Pagination
      emptyComponent={<NoResults />}
      hasNext={hasNext}
      hasPrev={hasPrev}
      isEmpty={isEmpty}
      isFetching={isFetching}
      isLoading={isLoading}
      loadingComponent={<Loading />}
      loadNext={loadNext}
      loadPrev={loadPrev}
    >
      <S.TestListContainer data-cy="test-list">
        {list?.map(environment => (
          <EnvironmentCard
            key={environment.name}
            environment={environment}
            setIsFormOpen={setIsFormOpen}
            setEnvironment={setEnvironment}
          />
        ))}
      </S.TestListContainer>
    </Pagination>
  );
};

export default EnvironmentList;
