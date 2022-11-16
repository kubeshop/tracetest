import {capitalize} from 'lodash';
import {Tooltip} from 'antd';
import {LinkOutlined} from '@ant-design/icons';
import {TestState} from 'constants/TestRun.constants';
import {TTest} from 'types/Test.types';
import {TTestRun, TTestRunState} from 'types/TestRun.types';
import TestRun from 'models/TestRun.model';
import * as S from './TransactionRunResult.styled';

const iconBasedOnResult = (result: TTestRunState, index: number) => {
  switch (result) {
    case TestState.FINISHED:
      return <S.IconSuccess />;
    case TestState.FAILED:
      return <S.IconSuccess />;
    default:
      return index + 1;
  }
};

interface IProps {
  index: number;
  test: TTest;
  testRun?: TTestRun;
}

const ExecutionStep = ({
  index,
  test: {name, trigger, id: testId},
  testRun: {id: runId, state, testVersion, passedAssertionCount, failedAssertionCount} = TestRun({}),
}: IProps) => {
  return (
    <S.Container data-cy={`run-card-${name}`} key={`${testId}-${runId}`}>
      <S.ExecutionStepStatus>{iconBasedOnResult(state, index)}</S.ExecutionStepStatus>
      <S.Info>
        <S.ExecutionStepName>{`${name} v${testVersion}`}</S.ExecutionStepName>
        <S.TagContainer>
          <S.TextTag>{trigger.method}</S.TextTag>
          <S.TextTag $isLight>{trigger.entryPoint}</S.TextTag>
          <S.TextTag>{capitalize(state)}</S.TextTag>
        </S.TagContainer>
      </S.Info>
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
        {runId && (
          <Tooltip title="Go to Run">
            <S.ExecutionStepRunLink to={`/test/${testId}/run/${runId}`} target="_blank">
              <LinkOutlined />
            </S.ExecutionStepRunLink>
          </Tooltip>
        )}
      </S.ExecutionStepStatus>
    </S.Container>
  );
};

export default ExecutionStep;
