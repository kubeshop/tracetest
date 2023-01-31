import EditTest from 'components/EditTest';
import FailedTrace from 'components/FailedTrace';
import RunDetailTriggerResponseFactory from 'components/RunDetailTriggerResponse/RunDetailTriggerResponseFactory';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import * as S from './RunDetailTrigger.styled';

interface IProps {
  test: Test;
  run: TestRun;
  isError: boolean;
}

const RunDetailTrigger = ({test, run: {state, triggerResult, triggerTime}, run, isError}: IProps) => {
  const shouldDisplayError = isError || state === TestState.FAILED;

  return (
    <S.Container>
      <S.SectionLeft>
        <EditTest test={test} />
      </S.SectionLeft>
      <S.SectionRight>
        {shouldDisplayError ? (
          <FailedTrace isDisplayingError={shouldDisplayError} run={run} />
        ) : (
          <RunDetailTriggerResponseFactory
            state={state}
            triggerResult={triggerResult}
            triggerTime={triggerTime}
            type={triggerResult?.type ?? TriggerTypes.http}
          />
        )}
      </S.SectionRight>
    </S.Container>
  );
};

export default RunDetailTrigger;
