/* eslint-disable */
import React, {useCallback, useEffect, useState} from 'react';
import styled from 'styled-components';

import {useStoreActions} from 'react-flow-renderer';
import {isEmpty} from 'lodash';

import {Button, Typography} from 'antd';
import {CloseCircleFilled} from '@ant-design/icons';

import 'react-reflex/styles.css';

import {useGetResultByIdQuery, useGetTestByIdQuery, useRunTestMutation} from 'redux/apis/Test.api';
import Diagram from 'components/Diagram';

import {GuidedTours} from 'services/GuidedTour.service';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Trace.styled';
import {ISpan} from '../../types/Span.types';
import {ITestRunResult} from '../../types/TestRunResult.types';
import {TestState} from '../../constants/TestRunResult.constants';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
import usePolling from '../../hooks/usePolling';
import {useAppDispatch} from '../../redux/hooks';
import {replace, updateTestResult} from '../../redux/slices/ResultList.slice';
import {SupportedDiagrams} from '../Diagram/Diagram';
import SpanDetail from '../SpanDetail';
import TraceTimeline from './TraceTimeline';
import {ResizableDrawer} from './ResizableDrawer';

const {onChangeTab} = TraceAnalyticsService;

const Grid = styled.div`
  display: flex;
  //height: calc(100vh - 200px);
  //overflow-y: scroll;
`;

type TraceProps = {
  testId: string;
  testResultId: string;
  onDismissTrace(): void;
  onRunTest(result: ITestRunResult): void;
};

const Trace: React.FC<TraceProps> = ({testId, testResultId, onDismissTrace, onRunTest}) => {
  const [selectedSpan, setSelectedSpan] = useState<ISpan | undefined>();
  const [isFirstLoad, setIsFirstLoad] = useState(true);
  const {data: test} = useGetTestByIdQuery(testId);
  const [runNewTest] = useRunTestMutation();
  const dispatch = useAppDispatch();

  const {
    data: testResultDetails,
    isError,
    refetch: refetchTrace,
  } = useGetResultByIdQuery({testId, resultId: testResultId});

  const addSelected = useStoreActions(actions => actions.addSelectedElements);

  const handleOnSpanSelected = useCallback(
    (spanId: string) => {
      addSelected([{id: spanId}]);
      const span = testResultDetails?.trace?.spans.find(({spanId: id}) => id === spanId);
      setSelectedSpan(span);
    },
    [addSelected, testResultDetails?.trace?.spans]
  );

  useGuidedTour(GuidedTours.Trace);
  usePolling({
    callback: refetchTrace,
    delay: 1000,
    isPolling:
      isError ||
      testResultDetails?.state === TestState.AWAITING_TRACE ||
      testResultDetails?.state === TestState.EXECUTING,
  });

  useEffect(() => {
    if (testResultDetails && test && !isFirstLoad) {
      dispatch(
        updateTestResult({
          trace: testResultDetails.trace!,
          resultId: testResultId,
          test,
        })
      );
    }
  }, [test, dispatch]);

  useEffect(() => {
    if (testResultDetails && !isEmpty(testResultDetails.trace) && !testResultDetails?.assertionResult && test) {
      setIsFirstLoad(false);
      dispatch(
        updateTestResult({
          trace: testResultDetails.trace!,
          resultId: testResultId,
          test,
        })
      );
    } else if (testResultDetails?.assertionResult && test) {
      setIsFirstLoad(false);

      dispatch(
        replace({
          resultId: testResultId,
          assertionResult: testResultDetails?.assertionResult!,
          test,
          trace: testResultDetails?.trace!,
        })
      );
    }
  }, [testResultDetails, test, testResultId, testId, dispatch]);

  const handleReRunTest = async () => {
    const result = await runNewTest(testId).unwrap();
    onRunTest(result);
  };

  if (isError || testResultDetails?.state === TestState.FAILED) {
    return (
      <S.FailedTrace>
        <CloseCircleFilled style={{color: 'red', fontSize: 32}} />
        <Typography.Title level={2}>Test Run Failed</Typography.Title>
        <div style={{display: 'grid', gap: 8, gridTemplateColumns: '1fr 1fr'}}>
          <Button onClick={handleReRunTest}>Rerun Test</Button>
          <Button onClick={onDismissTrace}>Cancel</Button>
        </div>
      </S.FailedTrace>
    );
  }

  const visiblePortion = 60;
  const height = `calc(100% - ${visiblePortion}px)`;
  const [max, setMax] = useState(600);
  return (
    <>
      <Grid style={{width: '100%', minHeight: height, maxHeight: height, height: height}}>
        <div style={{flexBasis: '50%', paddingTop: 10, paddingLeft: 10}}>
          <Diagram
            type={SupportedDiagrams.DAG}
            trace={testResultDetails?.trace!}
            onSelectSpan={handleOnSpanSelected}
            selectedSpan={selectedSpan}
          />
        </div>
        <div style={{flexBasis: '50%', overflowY: 'scroll', paddingTop: 10, paddingRight: 10}}>
          <SpanDetail resultId={testResultDetails?.resultId} testId={test?.testId} span={selectedSpan} />
        </div>
      </Grid>
      <ResizableDrawer min={visiblePortion} max={max}>
        <TraceTimeline
          min={visiblePortion}
          max={max}
          setMax={setMax}
          visiblePortion={visiblePortion}
          trace={testResultDetails?.trace!}
          onSelectSpan={handleOnSpanSelected}
          selectedSpan={selectedSpan}
        />
      </ResizableDrawer>
    </>
  );
};

export default Trace;
