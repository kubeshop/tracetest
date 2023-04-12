import {Button} from 'antd';
import RunActionsMenu from 'components/RunActionsMenu';
import TestActions from 'components/TestActions';
import TestState from 'components/TestState';
import {TestState as TestStateEnum} from 'constants/TestRun.constants';
import {isRunStateFinished} from 'models/TestRun.model';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import * as S from './RunDetailLayout.styled';

interface IProps {
  testId: string;
  testVersion: number;
}

const HeaderRight = ({testId, testVersion}: IProps) => {
  const {isDraftMode: isTestSpecsDraftMode} = useTestSpecs();
  const {isDraftMode: isTestOutputsDraftMode} = useTestOutput();
  const isDraftMode = isTestSpecsDraftMode || isTestOutputsDraftMode;
  const {run} = useTestRun();
  const {onRun} = useTest();
  const state = run.state;

  return (
    <S.Section $justifyContent="flex-end">
      {isDraftMode && <TestActions />}
      {!isDraftMode && state && state !== TestStateEnum.FINISHED && (
        <S.StateContainer data-cy="test-run-result-status">
          <S.StateText>Test status:</S.StateText>
          <TestState testState={state} />
        </S.StateContainer>
      )}
      {!isDraftMode && state && isRunStateFinished(state) && (
        <Button data-cy="run-test-button" ghost onClick={() => onRun(run.id)} type="primary">
          Run Test
        </Button>
      )}
      <RunActionsMenu
        isRunView
        resultId={run.id}
        testId={testId}
        testVersion={testVersion}
        transactionId={run.transactionId}
        transactionRunId={run.transactionRunId}
      />
    </S.Section>
  );
};

export default HeaderRight;
