import {CloseCircleOutlined} from '@ant-design/icons';
import {Button, Tooltip} from 'antd';
import RunActionsMenu from 'components/RunActionsMenu';
import TestActions from 'components/TestActions';
import TestState from 'components/TestState';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {isRunStateFinished, isRunStateStopped, isRunStateSucceeded} from 'models/TestRun.model';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import * as S from './RunDetailLayout.styled';
import EventLogPopover from '../EventLogPopover/EventLogPopover';
import RunStatusIcon from '../RunStatusIcon/RunStatusIcon';
import VariableSetSelector from '../VariableSetSelector/VariableSetSelector';

interface IProps {
  testId: string;
}

const HeaderRight = ({testId}: IProps) => {
  const {isDraftMode: isTestSpecsDraftMode} = useTestSpecs();
  const {isDraftMode: isTestOutputsDraftMode} = useTestOutput();
  const isDraftMode = isTestSpecsDraftMode || isTestOutputsDraftMode;
  const {
    isLoadingStop,
    run: {state, requiredGatesResult},
    run,
    stopRun,
    runEvents,
  } = useTestRun();
  const {onRun} = useTest();

  return (
    <S.Section $justifyContent="flex-end">
      {isDraftMode && <TestActions />}
      {!isDraftMode && state && state !== TestStateEnum.FINISHED && (
        <S.StateContainer data-cy="test-run-result-status">
          <S.StateText>Test status:</S.StateText>
          <TestState testState={state} />
          {state === TestStateEnum.AWAITING_TRACE && (
            <S.StopContainer>
              <Tooltip title="Stop test execution">
                <Button
                  disabled={isLoadingStop}
                  icon={<CloseCircleOutlined />}
                  onClick={stopRun}
                  shape="circle"
                  type="link"
                />
              </Tooltip>
            </S.StopContainer>
          )}
        </S.StateContainer>
      )}
      {(isRunStateSucceeded(state) || isRunStateStopped(state)) && (
        <RunStatusIcon state={state} requiredGatesResult={requiredGatesResult} />
      )}
      <VariableSetSelector />
      {!isDraftMode && state && isRunStateFinished(state) && (
        <Button data-cy="run-test-button" ghost onClick={() => onRun()} type="primary">
          Run Test
        </Button>
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
