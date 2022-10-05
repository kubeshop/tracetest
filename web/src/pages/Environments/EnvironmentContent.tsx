import {useCallback, useState} from 'react';
import SearchInput from '../../components/SearchInput';
import EnvironmentActions from './EnvironmentActions';
import {EnvironmentState} from './EnvironmentState';
import EnvironmentList from './EnvList';
import * as S from './Envs.styled';
import {EnvsModal} from './EnvsModal';
import {IEnvironment} from './IEnvironment';

const EnvironmentContent: React.FC = () => {
  const [state, setState] = useState<EnvironmentState>({query: '', dialog: false, environment: undefined});
  const onSearch = useCallback((value: string) => setState(st => ({...st, query: value})), [setState]);
  const openDialog = useCallback((dialog: boolean) => setState(st => ({...st, dialog})), [setState]);
  return (
    <S.Wrapper>
      <S.TitleText>All Envs</S.TitleText>
      <S.PageHeader>
        <SearchInput onSearch={onSearch} placeholder="Search test" />
        <EnvironmentActions openDialog={() => openDialog(true)} />
      </S.PageHeader>
      <EnvironmentList
        query={state.query}
        openDialog={openDialog}
        setEnvironment={(environment: IEnvironment) => setState(st => ({...st, environment}))}
      />
      <EnvsModal setState={setState} state={state} />
    </S.Wrapper>
  );
};

export default EnvironmentContent;
