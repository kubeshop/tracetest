import {useCallback, useState} from 'react';

import SearchInput from 'components/SearchInput';
import VariableSetsAnalytics from 'services/Analytics/VariableSetsAnalytics.service';
import VariableSet from 'models/VariableSet.model';
import {useVariableSet} from 'providers/VariableSet';
import * as S from './VariableSet.styled';
import VariableSetList from './VariableSetList';

const {onCreateVariableSetClick} = VariableSetsAnalytics;

const VariableSetContent = () => {
  const [query, setQuery] = useState<string>('');
  const onSearch = useCallback((value: string) => setQuery(value), [setQuery]);
  const {onOpenModal, onDelete} = useVariableSet();

  const handleOnClickCreate = useCallback(() => {
    onCreateVariableSetClick();
    onOpenModal();
  }, [onOpenModal]);

  const handleOnEdit = useCallback(
    (values: VariableSet) => {
      onOpenModal(values);
    },
    [onOpenModal]
  );

  return (
    <S.Wrapper>
      <S.MainHeaderContainer>
        <S.TitleText>All Variable Sets</S.TitleText>
      </S.MainHeaderContainer>

      <S.PageHeader>
        <SearchInput onSearch={onSearch} placeholder="Search variable set" />
        <S.ActionContainer>
          <S.CreateVarsButton onClick={handleOnClickCreate} type="primary">
            Create Variable Set
          </S.CreateVarsButton>
        </S.ActionContainer>
      </S.PageHeader>

      <VariableSetList onDelete={onDelete} onEdit={handleOnEdit} query={query} />
    </S.Wrapper>
  );
};

export default VariableSetContent;
