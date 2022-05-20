import {useCallback} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import TestHeader from 'components/TestHeader';
import {useGetResultByIdQueryPolling} from './hooks/useGetResultByIdQueryPolling';
import {TestState} from '../../constants/TestRun.constants';

import {useGetTestByIdQuery} from '../../redux/apis/TraceTest.api';
import useGuidedTour from '../../hooks/useGuidedTour';
import {visiblePortionFuction} from '../../utils/Common';
import {GuidedTours} from '../../services/GuidedTour.service';
import FailedTrace from '../../components/FailedTrace';
import Trace from '../../components/Trace';
import * as S from './Trace.styled';
import {useTestRun} from '../../providers/TestRun/TestRun.provider';

const TraceContent = () => {
  const {testId = ''} = useParams();
  const {visiblePortion, height} = visiblePortionFuction();
  const navigate = useNavigate();
  const {data: test} = useGetTestByIdQuery({testId});
  useGuidedTour(GuidedTours.Trace);

  const {isError, run, refetch} = useTestRun();
  const isDisplayingError = isError || run.state === TestState.FAILED;

  useGetResultByIdQueryPolling(refetch, isError, run);

  const onRunTest = useCallback(() => {
    console.log('onRunTest');
  }, []);

  return test ? (
    <S.Wrapper>
      <TestHeader test={test} onBack={() => navigate(`/test/${testId}`)} testState={run.state} />
      <FailedTrace
        onRunTest={onRunTest}
        testId={testId}
        isDisplayingError={isDisplayingError}
        onEdit={() => console.log('onEdit')}
      />
      <Trace
        displayError={isDisplayingError}
        minHeight={height}
        run={run}
        test={test}
        visiblePortion={visiblePortion}
      />
    </S.Wrapper>
  ) : null;
};

export default TraceContent;
