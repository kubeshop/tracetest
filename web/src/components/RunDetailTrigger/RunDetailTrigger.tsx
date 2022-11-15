import EditTest from 'components/EditTest';
import {TestState} from '../../constants/TestRun.constants';
import {TTest} from '../../types/Test.types';
import {TTestRun} from '../../types/TestRun.types';
import FailedTrace from '../FailedTrace';
import RunDetailTriggerResponse from '../RunDetailTriggerResponse';
import * as S from './RunDetailTrigger.styled';

interface IProps {
  test: TTest;
  run: TTestRun;
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
          <RunDetailTriggerResponse triggerResult={triggerResult} triggerTime={triggerTime} />
        )}
      </S.SectionRight>
    </S.Container>
  );
};

export default RunDetailTrigger;
