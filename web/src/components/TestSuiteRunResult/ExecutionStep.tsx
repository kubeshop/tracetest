import {LinkOutlined} from '@ant-design/icons';
import {Tooltip} from 'antd';
import {capitalize} from 'lodash';
import {useState} from 'react';
import KeyValueRow from 'components/KeyValueRow';
import Link from 'components/Link';
import {TestState} from 'constants/TestRun.constants';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import * as S from './TestSuiteRunResult.styled';
import RunStatusIcon from '../RunStatusIcon';

interface IProps {
  test: Test;
  testRun?: TestRun;
  hasRunFailed: boolean;
}

const ExecutionStep = ({
  test: {name, trigger, id: testId},
  hasRunFailed,
  testRun: {
    id: runId,
    state,
    testVersion,
    passedAssertionCount,
    failedAssertionCount,
    outputs,
    requiredGatesResult,
  } = TestRun({
    state: hasRunFailed ? TestState.SKIPPED : TestState.WAITING,
  }),
}: IProps) => {
  const [toggleOutputs, setToggleOutputs] = useState(false);
  const stateIsFinished = ([TestState.FINISHED, TestState.FAILED] as string[]).includes(state);
  const toLink = runId ? `/test/${testId}/run/${runId}` : `/test/${testId}`;

  return (
    <S.Container data-cy={`testsuite-execution-step-${name}`}>
      <S.Content>
        <RunStatusIcon state={state} requiredGatesResult={requiredGatesResult} />
        <Link to={toLink} target="_blank">
          <S.Info>
            <S.ItemName>{`${name} v${testVersion}`}</S.ItemName>
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
            <div>
              <S.ExecutionStepRunLink to={toLink} target="_blank" data-cy="execution-step-run-link">
                <LinkOutlined />
              </S.ExecutionStepRunLink>
            </div>
          </Tooltip>
        </S.ExecutionStepStatus>
      </S.Content>

      <S.OutputsContainer>
        {!!outputs?.length && (
          <S.OutputsButton onClick={() => setToggleOutputs(prev => !prev)} type="link">
            {toggleOutputs ? 'Hide Outputs' : 'Show Outputs'}
          </S.OutputsButton>
        )}

        {toggleOutputs && (
          <S.OutputsContent>
            {outputs?.map?.(output => (
              <KeyValueRow key={output.name} keyName={output.name} value={output.value} />
            ))}
          </S.OutputsContent>
        )}
      </S.OutputsContainer>
    </S.Container>
  );
};

export default ExecutionStep;
