import {useTheme} from 'styled-components';
import TestRun, {isRunStateFailed, isRunStateFinished, isRunStateStopped} from 'models/TestRun.model';
import RequiredGatesResult from 'models/RequiredGatesResult.model';
import {ToTitle} from 'utils/Common';
import * as S from './RunStatusIcon.styled';

interface IProps {
  state: TestRun['state'];
  requiredGatesResult: RequiredGatesResult;
}

const RunTooltipTitle = ({requiredGatesResult: {required, requiredFailedGates, passed}, state}: IProps) => {
  const {
    color: {error, success},
  } = useTheme();

  const isStopped = isRunStateStopped(state);
  if (isStopped) return <>The run has been manually stopped</>;

  const isFinished = isRunStateFinished(state);
  const isSuccessful = isFinished && passed;
  if (isSuccessful) {
    return (
      <>
        {required.length ? 'Required Gates' : 'No gates were required'}
        <S.GateList $color={success}>
          {required.map(gate => (
            <S.Gate key={gate}>
              <S.GateName>{ToTitle(gate)}</S.GateName>
            </S.Gate>
          ))}
        </S.GateList>
      </>
    );
  }

  const isRunnerFailed = isRunStateFailed(state);
  if (isRunnerFailed) return <>The run execution failed</>;

  const isGatesFailed = isFinished && !passed;
  if (isGatesFailed)
    return (
      <>
        Failed Required Gates
        <S.GateList $color={error}>
          {requiredFailedGates.map(gate => (
            <S.Gate key={gate}>
              <S.GateName>{ToTitle(gate)}</S.GateName>
            </S.Gate>
          ))}
        </S.GateList>
      </>
    );

  return <>The run execution is in progress</>;
};

export default RunTooltipTitle;
