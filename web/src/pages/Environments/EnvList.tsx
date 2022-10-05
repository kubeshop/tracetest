import {useGetEnvListQuery} from '../../redux/apis/TraceTest.api';
import {EnvironmentCard} from './EnvironmentCard';
import * as S from './Envs.styled';
import {IEnvironment} from './IEnvironment';

interface IProps {
  query: string;
  openDialog: (mode: boolean) => void;
  setEnvironment: (mode: IEnvironment) => void;
}

const EnvList = ({query, openDialog, setEnvironment}: IProps) => {
  const {data: list} = useGetEnvListQuery({query});
  return (
    <S.TestListContainer data-cy="test-list">
      {list?.map((environment: IEnvironment) => (
        <EnvironmentCard
          key={environment.name}
          environment={environment}
          openDialog={openDialog}
          setEnvironment={setEnvironment}
        />
      ))}
    </S.TestListContainer>
  );
};

export default EnvList;
