import {capitalize} from 'lodash';
import {Tooltip} from 'antd';
import {LinkOutlined} from '@ant-design/icons';
import {Link} from 'react-router-dom';
import {TestState} from 'constants/TestRun.constants';
import {TTest} from 'types/Test.types';
import {TTestRun, TTestRunState} from 'types/TestRun.types';
import TestRun from 'models/TestRun.model';
import * as S from './TransactionRunResult.styled';

const iconBasedOnResult = (state: TTestRunState, failedAssertions: number, index: number) => {
  if (state !== TestState.FAILED && state !== TestState.FINISHED) {
    return null;
  }

  if (state === TestState.FAILED || failedAssertions > 0) {
    return <S.IconFail />;
  }
  if (state === TestState.FINISHED || failedAssertions === 0) {
    return <S.IconSuccess />;
  }

  return index + 1;
};

interface IProps {
  index: number;
  test: TTest;
  testRun?: TTestRun;
  hasRunFailed: boolean;
}

const ExecutionStep = ({
  index,
  test: {name, trigger, id: testId},
  hasRunFailed,
  testRun: {id: runId, state, testVersion, passedAssertionCount, failedAssertionCount} = TestRun({
    state: hasRunFailed ? TestState.SKIPPED : TestState.WAITING,
  }),
}: IProps) => {
  const stateIsFinished = ([TestState.FINISHED, TestState.FAILED] as string[]).includes(state);
  const toLink = runId ? `/test/${testId}/run/${runId}` : `/test/${testId}`;

  return (
    <S.Container data-cy={`transaction-execution-step-${name}`}>
      <S.ExecutionStepStatus>{iconBasedOnResult(state, failedAssertionCount, index)}</S.ExecutionStepStatus>
      <Link to={toLink} target="_blank">
        <S.Info>
          <S.ExecutionStepName>{`${name} v${testVersion}`}</S.ExecutionStepName>
          <S.TagContainer>
            <S.TextTag>{trigger.method}</S.TextTag>
            <S.EntryPointTag $isLight>{trigger.entryPoint}</S.EntryPointTag>
            {!stateIsFinished && <S.TextTag>{capitalize(state)}</S.TextTag>}
          </S.TagContainer>
        </S.Info>
      </Link>
      <S.AssertionResultContainer>
        {runId && (
          <>
            <Tooltip title="Passed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed />
                {passedAssertionCount}
              </S.HeaderDetail>
            </Tooltip>
            <Tooltip title="Failed assertions">
              <S.HeaderDetail>
                <S.HeaderDot $passed={false} />
                {failedAssertionCount}
              </S.HeaderDetail>
            </Tooltip>
          </>
        )}
      </S.AssertionResultContainer>
      <S.ExecutionStepStatus>
        <Tooltip title="Go to Run">
          <S.ExecutionStepRunLink to={toLink} target="_blank" data-cy="execution-step-run-link">
            <LinkOutlined />
          </S.ExecutionStepRunLink>
        </Tooltip>
      </S.ExecutionStepStatus>
    </S.Container>
  );
};

export default ExecutionStep;
