import {Tooltip} from 'antd';
import TestRun, {isRunStateFailed, isRunStateFinished, isRunStateStopped} from 'models/TestRun.model';
import RequiredGatesResult from 'models/RequiredGatesResult.model';
import * as S from './RunStatusIcon.styled';
import RunTooltipTitle from './RunTooltipTitle';

interface IProps {
  state: TestRun['state'];
  requiredGatesResult: RequiredGatesResult;
}

const Icon = ({state, requiredGatesResult: {passed}}: IProps) => {
  if (!isRunStateFinished(state)) {
    return <S.LoadingIcon />;
  }

  if (isRunStateStopped(state)) {
    return <S.IconInfo />;
  }

  if (isRunStateFailed(state) || !passed) {
    return <S.IconFail />;
  }

  return <S.IconSuccess />;
};

const RunStatusIcon = (props: IProps) => {
  return (
    <Tooltip title={() => <RunTooltipTitle {...props} />}>
      <S.IconWrapper>
        <Icon {...props} />
      </S.IconWrapper>
    </Tooltip>
  );
};

export default RunStatusIcon;
