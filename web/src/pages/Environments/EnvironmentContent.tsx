import React, {useState} from 'react';
import Pagination from '../../components/Pagination';
import SearchInput from '../../components/SearchInput';
import {useGetEnvListQuery} from '../../redux/apis/TraceTest.api';
import * as S from './Environment.styled';
import EnvironmentActions from './EnvironmentActions';
import {EnvironmentCard} from './EnvironmentCard';
import {EnvironmentModal} from './EnvironmentModal';
import {IEnvironment} from './IEnvironment';

const EnvironmentContent: React.FC = () => {
  const [isFormOpen, setIsFormOpen] = useState<boolean>(false);
  const [environment, setEnvironment] = useState<IEnvironment | undefined>(undefined);
  return (
    <S.Wrapper>
      <S.TitleText>All Envs</S.TitleText>
      <Pagination<IEnvironment, {query: string}> query={useGetEnvListQuery} defaultParameters={{query: ''}}>
        {(pagination, [, setParams]) => {
          const onSearch = (value: string) => setParams(st => ({...st, query: value}));
          return (
            <>
              <S.PageHeader>
                <SearchInput onSearch={onSearch} placeholder="Search test" />
                <EnvironmentActions setIsFormOpen={setIsFormOpen} />
              </S.PageHeader>
              <S.TestListContainer data-cy="test-list">
                {pagination.list?.map(env => (
                  <EnvironmentCard
                    key={env.name}
                    environment={env}
                    setIsFormOpen={setIsFormOpen}
                    setEnvironment={setEnvironment}
                  />
                ))}
              </S.TestListContainer>
            </>
          );
        }}
      </Pagination>
      <EnvironmentModal
        isFormOpen={isFormOpen}
        environment={environment}
        setEnvironment={setEnvironment}
        setIsFormOpen={setIsFormOpen}
      />
    </S.Wrapper>
  );
};

export default EnvironmentContent;
