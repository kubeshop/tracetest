import TransactionRun from 'models/TransactionRun.model';
import {TestState} from 'constants/TestRun.constants';
import * as S from './RunStatusIcon.styled';

interface IProps {
  state: TransactionRun['state'];
  hasFailedTests: boolean;
}

const Icon = ({state, hasFailedTests}: IProps) => {
  if (state !== TestState.FAILED && state !== TestState.FINISHED) {
    return <S.LoadingIcon />;
  }

  if (state === TestState.FAILED || hasFailedTests) {
    return <S.IconFail />;
  }

  return <S.IconSuccess />;
};

const TransactionRunStatusIcon = (props: IProps) => {
  return (
    <S.IconWrapper>
      <Icon {...props} />
    </S.IconWrapper>
  );
};

export default TransactionRunStatusIcon;
