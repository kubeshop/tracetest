import {useCallback, useState} from 'react';
import {withTracker} from 'ga-4-react';
import {useNavigate, useParams} from 'react-router-dom';
import {ReactFlowProvider} from 'react-flow-renderer';
import Layout from 'components/Layout';
import TestHeader from 'components/TestHeader';
import {useGetResultByIdQueryPolling} from './hooks/useGetResultByIdQueryPolling';
import {useUpdateTestResultEffect} from './hooks/useUpdateTestResultEffect';
import {useUpdateFirstLoadEffect} from './hooks/useUpdateFirstLoadEffect';
import {TestState} from '../../constants/TestRunResult.constants';
import AssertionFormProvider from '../../components/AssertionForm/AssertionFormProvider';

import {useGetResultByIdQuery, useGetTestByIdQuery} from '../../redux/apis/Test.api';
import {useAppDispatch} from '../../redux/hooks';
import useGuidedTour from '../../hooks/useGuidedTour';
import {visiblePortionFuction} from '../../utils/Common';
import {GuidedTours} from '../../services/GuidedTour.service';
import FailedTrace from '../../components/FailedTrace';
import Trace from '../../components/Trace';
import * as S from './Trace.styled';

const TracePage = () => {
  const dispatch = useAppDispatch();
  const {testId = '', resultId = ''} = useParams();
  const {visiblePortion, height} = visiblePortionFuction();
  const [isFirstLoad, setIsFirstLoad] = useState(true);
  const navigate = useNavigate();
  const {data: test} = useGetTestByIdQuery(testId);
  useGuidedTour(GuidedTours.Trace);

  const {isError, data: result, refetch} = useGetResultByIdQuery({testId, resultId});
  const isDisplayingError = isError || result?.state === TestState.FAILED;

  useGetResultByIdQueryPolling(refetch, isError, result);
  useUpdateTestResultEffect(result, test, isFirstLoad, dispatch, resultId);
  useUpdateFirstLoadEffect(result, test, setIsFirstLoad, dispatch, resultId, testId);

  const onRunTest = useCallback(() => {
    console.log('onRunTest');
  }, []);

  return test && result ? (
    <ReactFlowProvider>
      <AssertionFormProvider testId={testId}>
        <Layout>
          <S.Wrapper>
            <TestHeader test={test} onBack={() => navigate(`/test/${testId}`)} testState={result.state} />
            <FailedTrace
              onRunTest={onRunTest}
              testId={testId}
              isDisplayingError={isDisplayingError}
              onEdit={() => console.log('onEdit')}
            />
            <Trace
              displayError={isDisplayingError}
              minHeight={height}
              testResultDetails={result}
              test={test}
              visiblePortion={visiblePortion}
            />
          </S.Wrapper>
        </Layout>
      </AssertionFormProvider>
    </ReactFlowProvider>
  ) : null;
};

export default withTracker(TracePage);
