import {useCallback, useEffect, useState} from 'react';
import styled from 'styled-components';

import {useStoreActions} from 'react-flow-renderer';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
import {isEmpty} from 'lodash';

import {Button, Tabs, Typography} from 'antd';
import {CloseCircleFilled} from '@ant-design/icons';

import 'react-reflex/styles.css';

import {useGetTestByIdQuery, useGetResultByIdQuery, useRunTestMutation} from 'redux/apis/Test.api';
import Diagram from 'components/Diagram';

import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {Steps} from 'components/GuidedTour/traceStepList';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Trace.styled';

import SpanDetail from '../SpanDetail';
import TestResults from './TestResults';
import {ISpan} from '../../types/Span.types';
import {ITestRunResult} from '../../types/TestRunResult.types';
import {TestState} from '../../constants/TestRunResult.constants';
import TraceTimeline from './TraceTimeline';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
import usePolling from '../../hooks/usePolling';
import {useAppDispatch} from '../../redux/hooks';
import {replace, updateTestResult} from '../../redux/slices/ResultList.slice';
import { SupportedDiagrams } from '../Diagram/Diagram';

const {onChangeTab} = TraceAnalyticsService;

const Grid = styled.div`
  display: grid;
  height: calc(100vh - 200px);
  overflow: scroll;
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

  return (
    <main>
      <Grid>
        <ReflexContainer style={{height: '100vh'}} orientation="horizontal">
          <ReflexElement flex={0.6}>
            <ReflexContainer orientation="vertical">
              <ReflexElement
                flex={0.5}
                className="left-pane"
                data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Diagram)}
              >
                <div className="pane-content">
                  <Diagram
                    type={SupportedDiagrams.DAG}
                    trace={testResultDetails?.trace!}
                    onSelectSpan={handleOnSpanSelected}
                    selectedSpan={selectedSpan}
                  />
                </div>
              </ReflexElement>
              <ReflexElement flex={0.5} className="right-pane">
                <div className="pane-content" style={{padding: '14px 24px', overflow: 'hidden'}}>
                  <S.TraceTabs onChange={activeTab => onChangeTab(activeTab)}>
                    <Tabs.TabPane
                      tab={
                        <span data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.SpanDetail)}>
                          Span Detail
                        </span>
                      }
                      key="span-detail"
                    >
                      <SpanDetail
                        resultId={testResultDetails?.resultId}
                        testId={test?.testId}
                        span={selectedSpan}
                      />
                    </Tabs.TabPane>
                    <Tabs.TabPane
                      tab={
                        <span data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.TestResults)}>
                          Test Results
                        </span>
                      }
                      key="test-results"
                    >
                      <TestResults
                        onSpanSelected={handleOnSpanSelected}
                        trace={testResultDetails?.trace}
                        resultId={testResultId}
                      />
                    </Tabs.TabPane>
                  </S.TraceTabs>
                </div>
              </ReflexElement>
            </ReflexContainer>
          </ReflexElement>
          <ReflexSplitter />
          <ReflexElement>
            <div className="pane-content">
              <TraceTimeline
                trace={testResultDetails?.trace!}
                onSelectSpan={handleOnSpanSelected}
                selectedSpan={selectedSpan}
              />
            </div>
          </ReflexElement>
        </ReflexContainer>
      </Grid>
    </main>
  );
};

export default Trace;
