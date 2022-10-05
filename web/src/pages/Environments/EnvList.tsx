import Pagination from '../../components/Pagination';
import usePagination from '../../hooks/usePagination';
import {useGetEnvListQuery} from '../../redux/apis/TraceTest.api';
import Loading from '../Home/Loading';
import NoResults from '../Home/NoResults';
import {EnvironmentCard} from './EnvironmentCard';
import * as S from './Envs.styled';
import {IEnvironment} from './IEnvironment';

interface IProps {
  query: string;
  openDialog: (mode: boolean) => void;
  setEnvironment: (mode: IEnvironment) => void;
}

const EnvList = ({query, openDialog, setEnvironment}: IProps) => {
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
            openDialog={openDialog}
            setEnvironment={setEnvironment}
          />
        ))}
      </S.TestListContainer>
    </Pagination>
  );
};

export default EnvList;
