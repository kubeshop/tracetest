import {TestState} from 'constants/TestRun.constants';
import TestRun, {isRunStateFailed, isRunStateFinished, isRunStateStopped} from 'models/TestRun.model';
import RequiredGatesResult from 'models/RequiredGatesResult.model';
import {useTheme} from 'styled-components';
import {TTestRunState} from 'types/TestRun.types';
import {ToTitle} from 'utils/Common';
import * as S from './RunStatusIcon.styled';

function getRunStateFailedMessage(state: TTestRunState) {
  switch (state) {
    case TestState.TRIGGER_FAILED:
      return 'The run failed in the trigger stage';
    case TestState.TRACE_FAILED:
      return 'The run failed in fetching the trace';
    case TestState.ASSERTION_FAILED:
      return 'The run failed to execute the assertions';
    default:
      return 'The run execution failed';
  }
}

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

  const isFailed = isRunStateFailed(state);
  if (isFailed) return <>{getRunStateFailedMessage(state)}</>;

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
