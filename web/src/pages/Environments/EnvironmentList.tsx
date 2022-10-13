import {Dispatch, SetStateAction} from 'react';
import Pagination from '../../components/Pagination';
import {useGetEnvListQuery} from '../../redux/apis/TraceTest.api';
import * as S from './Environment.styled';
import {EnvironmentCard} from './EnvironmentCard';
import {IEnvironment} from './IEnvironment';

interface IProps {
  query: string;
  setIsFormOpen: Dispatch<SetStateAction<boolean>>;
  setEnvironment: Dispatch<SetStateAction<IEnvironment | undefined>>;
}

const EnvironmentList = ({query, setEnvironment, setIsFormOpen}: IProps) => {
  return (
    <Pagination<IEnvironment, {query: string}> query={useGetEnvListQuery} defaultParameters={{query}}>
      {pagination => (
        <S.TestListContainer data-cy="test-list">
          {pagination.list?.map(environment => (
            <EnvironmentCard
              key={environment.name}
              environment={environment}
              setIsFormOpen={setIsFormOpen}
              setEnvironment={setEnvironment}
            />
          ))}
        </S.TestListContainer>
      )}
    </Pagination>
  );
};

export default EnvironmentList;
