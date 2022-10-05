import {Dispatch, SetStateAction} from 'react';
import EnvironmentsAnalytics from 'services/Analytics/EnvironmentsAnalytics.service';
import * as S from './Environment.styled';

const {onCreateEnvironmentClick} = EnvironmentsAnalytics;

interface IProps {
  setIsFormOpen: Dispatch<SetStateAction<boolean>>;
}

const EnvironmentActions = ({setIsFormOpen}: IProps): React.ReactElement => {
  return (
    <S.ActionContainer>
      <S.CreateEnvironmentButton
        data-cy="create-test-button"
        type="primary"
        onClick={() => {
          setIsFormOpen(true);
          onCreateEnvironmentClick();
        }}
      >
        Create Environment
      </S.CreateEnvironmentButton>
    </S.ActionContainer>
  );
};

export default EnvironmentActions;
