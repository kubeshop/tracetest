import CreateButton from 'components/CreateButton';
import RunActionsMenu from 'components/RunActionsMenu';
import TestActions from 'components/TestActions';
import TestState from 'components/TestState';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import Test from 'models/Test.model';
import {isRunPollingState, isRunStateFinished, isRunStateStopped, isRunStateSucceeded} from 'models/TestRun.model';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import * as S from './RunDetailLayout.styled';
import EventLogPopover from '../EventLogPopover/EventLogPopover';
import RunStatusIcon from '../RunStatusIcon/RunStatusIcon';
import VariableSetSelector from '../VariableSetSelector/VariableSetSelector';
import TracePollingActions from '../SkipPollingPopover/SkipPollingPopover';
import useSkipPolling from './hooks/useSkipPolling';

interface IProps {
  testId: string;
  triggerType: TriggerTypes;
}

const HeaderRight = ({testId, triggerType}: IProps) => {
  const {isDraftMode: isTestSpecsDraftMode} = useTestSpecs();
  const {isDraftMode: isTestOutputsDraftMode} = useTestOutput();
  const isDraftMode = isTestSpecsDraftMode || isTestOutputsDraftMode;
  const {
    run: {state, requiredGatesResult, createdAt},
    run,
    runEvents,
  } = useTestRun();
  const {onRun} = useTest();

  const {onSkipPolling, isLoading} = useSkipPolling();

  return (
    <S.Section $justifyContent="flex-end">
      {isDraftMode && <TestActions />}
      {!isDraftMode && state && state !== TestStateEnum.FINISHED && (
        <S.StateContainer data-cy="test-run-result-status">
          <S.StateText>Test status:</S.StateText>
          <TestState testState={state} />
          {isRunPollingState(state) && (
            <TracePollingActions startTime={createdAt} isLoading={isLoading} skipPolling={onSkipPolling} />
          )}
        </S.StateContainer>
      )}
      {(isRunStateSucceeded(state) || isRunStateStopped(state)) && (
        <RunStatusIcon state={state} requiredGatesResult={requiredGatesResult} />
      )}
      <VariableSetSelector />
      {!isDraftMode && state && isRunStateFinished(state) && Test.shouldAllowRun(triggerType) && (
        <CreateButton data-cy="run-test-button" ghost onClick={() => onRun()} type="primary">
          Run Test
        </CreateButton>
      )}
      <EventLogPopover runEvents={runEvents} />
      <RunActionsMenu
        isRunView
        resultId={run.id}
        testId={testId}
        testSuiteRunId={run.testSuiteRunId}
        testSuiteId={run.testSuiteId}
      />
    </S.Section>
  );
};

export default HeaderRight;
