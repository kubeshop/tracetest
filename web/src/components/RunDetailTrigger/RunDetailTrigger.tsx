import EditTest from 'components/EditTest';
import RunDetailTriggerResponseFactory from 'components/RunDetailTriggerResponse/RunDetailTriggerResponseFactory';
import RunEvents from 'components/RunEvents';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunDetailTrigger.styled';

interface IProps {
  test: Test;
  run: TestRun;
  runEvents: TestRunEvent[];
  isError: boolean;
}

const RunDetailTrigger = ({test, run: {id, state, triggerResult, triggerTime}, runEvents, isError}: IProps) => {
  const shouldDisplayError = isError || state === TestState.TRIGGER_FAILED;

  return (
    <S.Container>
      <S.SectionLeft>
        <EditTest test={test} />
      </S.SectionLeft>
      <S.SectionRight>
        {shouldDisplayError ? (
          <RunEvents events={runEvents} stage={TestRunStage.Trigger} state={state} />
        ) : (
          <RunDetailTriggerResponseFactory
            runId={id}
            state={state}
            testId={test.id}
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
