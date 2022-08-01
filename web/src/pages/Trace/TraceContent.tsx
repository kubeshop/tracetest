import {useNavigate, useParams} from 'react-router-dom';

import FailedTrace from 'components/FailedTrace';
import Run from 'components/Run';
import TestHeader from 'components/TestHeader';
import TraceActions from 'components/TraceActions';
import {TestState} from 'constants/TestRun.constants';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './Trace.styled';

const TraceContent = () => {
  const {testId = ''} = useParams();
  const navigate = useNavigate();
  const {isDraftMode, test} = useTestDefinition();

  const {isError, run} = useTestRun();
  const isDisplayingError = isError || run.state === TestState.FAILED;

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
      <FailedTrace testId={testId} run={run} isDisplayingError={isDisplayingError} />
      <Run displayError={isDisplayingError} run={run} test={test} />
    </S.Wrapper>
  ) : null;
};

export default TraceContent;
