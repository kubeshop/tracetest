import EnvironmentsAnalytics from 'services/Analytics/EnvironmentsAnalytics.service';
import * as S from './Envs.styled';

const {onCreateEnvironmentClick} = EnvironmentsAnalytics;

interface IProps {
  openDialog: () => void;
}

const EnvironmentActions = ({openDialog}: IProps): React.ReactElement => {
  return (
    <S.ActionContainer>
      <S.CreateEnvironmentButton
        data-cy="create-test-button"
        type="primary"
        onClick={() => {
          openDialog();
          onCreateEnvironmentClick();
        }}
      >
        Create Environment
      </S.CreateEnvironmentButton>
    </S.ActionContainer>
  );
};

export default EnvironmentActions;
