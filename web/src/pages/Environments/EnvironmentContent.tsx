import {useCallback, useState} from 'react';

import SearchInput from 'components/SearchInput';
import EnvironmentsAnalytics from 'services/Analytics/EnvironmentsAnalytics.service';
import Environment from 'models/Environment.model';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import * as S from './Environment.styled';
import EnvironmentList from './EnvironmentList';

const {onCreateEnvironmentClick} = EnvironmentsAnalytics;

const EnvironmentContent = () => {
  const [query, setQuery] = useState<string>('');
  const onSearch = useCallback((value: string) => setQuery(value), [setQuery]);
  const {onOpenModal, onDelete} = useEnvironment();

  const handleOnClickCreate = useCallback(() => {
    onCreateEnvironmentClick();
    onOpenModal();
  }, [onOpenModal]);

  const handleOnEdit = useCallback(
    (values: Environment) => {
      onOpenModal(values);
    },
    [onOpenModal]
  );

  return (
    <S.Wrapper>
      <S.MainHeaderContainer>
        <S.TitleText>All Environments</S.TitleText>
      </S.MainHeaderContainer>

      <S.PageHeader>
        <SearchInput onSearch={onSearch} placeholder="Search environment" />
        <S.ActionContainer>
          <S.CreateEnvironmentButton onClick={handleOnClickCreate} type="primary">
            Create Environment
          </S.CreateEnvironmentButton>
        </S.ActionContainer>
      </S.PageHeader>

      <EnvironmentList onDelete={onDelete} onEdit={handleOnEdit} query={query} />
    </S.Wrapper>
  );
};

export default EnvironmentContent;
