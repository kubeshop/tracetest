import React, {useCallback, useState} from 'react';
import SearchInput from '../../components/SearchInput';
import * as S from './Environment.styled';
import EnvironmentActions from './EnvironmentActions';
import EnvironmentList from './EnvironmentList';
import {EnvironmentModal} from './EnvironmentModal';
import {TEnvironment} from '../../types/Environment.types';

const EnvironmentContent: React.FC = () => {
  const [query, setQuery] = useState<string>('');
  const [isFormOpen, setIsFormOpen] = useState<boolean>(false);
  const [environment, setEnvironment] = useState<TEnvironment | undefined>(undefined);
  const onSearch = useCallback((value: string) => setQuery(value), [setQuery]);
  return (
    <S.Wrapper>
      <S.TitleText>All Envs</S.TitleText>
      <S.PageHeader>
        <SearchInput onSearch={onSearch} placeholder="Search test" />
        <EnvironmentActions setIsFormOpen={setIsFormOpen} />
      </S.PageHeader>
      <EnvironmentList query={query} setIsFormOpen={setIsFormOpen} setEnvironment={setEnvironment} />
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
