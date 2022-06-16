import {useCallback} from 'react';
import {useNavigate, useParams} from 'react-router-dom';

import FailedTrace from 'components/FailedTrace';
import Run from 'components/Run';
import TestHeader from 'components/TestHeader';
import TraceActions from 'components/TraceActions';
import {TestState} from 'constants/TestRun.constants';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useGetResultByIdQueryPolling} from './hooks/useGetResultByIdQueryPolling';
import * as S from './Trace.styled';

const TraceContent = () => {
  const {testId = ''} = useParams();
  const navigate = useNavigate();
  const {isDraftMode, test} = useTestDefinition();

  const {isError, run, refetch} = useTestRun();
  const isDisplayingError = isError || run.state === TestState.FAILED;

  useGetResultByIdQueryPolling(refetch, isError, run);

  const onRunTest = useCallback(() => {
    console.log('onRunTest');
  }, []);

  return test ? (
    <S.Wrapper>
      <TestHeader
        executionTime={run?.executionTime}
        extraContent={isDraftMode ? <TraceActions /> : undefined}
        onBack={() => navigate(`/test/${testId}`)}
        showInfo
        test={test}
        testState={run.state}
        testVersion={run.testVersion}
        totalSpans={run?.trace?.spans?.length}
      />
      <FailedTrace
        onRunTest={onRunTest}
        testId={testId}
        run={run}
        isDisplayingError={isDisplayingError}
        onEdit={() => console.log('onEdit')}
      />
      <Run displayError={isDisplayingError} run={run} test={test} />
    </S.Wrapper>
  ) : null;
};

export default TraceContent;
