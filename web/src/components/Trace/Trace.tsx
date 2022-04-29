import {useCallback, useEffect, useMemo, useState} from 'react';
import styled from 'styled-components';

import {useStoreActions} from 'react-flow-renderer';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
import {isEmpty} from 'lodash';

import {Button, Tabs, Typography} from 'antd';
import {CloseCircleFilled} from '@ant-design/icons';

import 'react-reflex/styles.css';

import {
  useGetTestByIdQuery,
  useGetTestResultByIdQuery,
  useRunTestMutation,
  useUpdateTestResultMutation,
} from 'gateways/Test.gateway';
import {
  parseAssertionResultListToTestResult,
  parseTestResultToAssertionResultList,
  runTest,
} from 'services/Trace.service';
import TraceDiagram from 'components/TraceDiagram/TraceDiagram';

import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {Steps} from 'components/GuidedTour/traceStepList';
import useGuidedTour from 'hooks/useGuidedTour';
import * as S from './Trace.styled';

import SpanDetail from './SpanDetail';
import TestResults from './TestResults';
import {ISpan} from '../../types/Span.types';
import {IAssertionResult} from '../../types/Assertion.types';
import {ITestRunResult} from '../../types/TestRunResult.types';
import {TestState} from '../../constants/TestRunResult.constants';
import TraceTimeline from './TraceTimeline';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';

const {onChangeTab} = TraceAnalyticsService;

const Grid = styled.div`
  display: grid;
  height: calc(100vh - 200px);
  overflow: scroll;
`;

export type TSpanInfo = {
  id: string;
  parentIds: string[];
  data: ISpan;
};

export type TSpanMap = Record<string, TSpanInfo>;

type TraceProps = {
  testId: string;
  testResultId: string;
  onDismissTrace(): void;
  onRunTest(result: ITestRunResult): void;
};

const Trace: React.FC<TraceProps> = ({testId, testResultId, onDismissTrace, onRunTest}) => {
  const [selectedSpan, setSelectedSpan] = useState<TSpanInfo | undefined>();
  const [traceResultList, setTraceResultList] = useState<IAssertionResult[]>([]);
  const [isFirstLoad, setIsFirstLoad] = useState(true);
  const [updateTestResult] = useUpdateTestResultMutation();
  const {data: test} = useGetTestByIdQuery(testId);
  const [runNewTest] = useRunTestMutation();

  const {
    data: testResultDetails,
    isError,
    refetch: refetchTrace,
  } = useGetTestResultByIdQuery({testId, resultId: testResultId});

  const spanMap = useMemo<TSpanMap>(() => {
    return testResultDetails?.trace?.resourceSpans
      ?.map(i => i.instrumentationLibrarySpans.map((el: any) => el.spans))
      ?.flat(2)
      ?.reduce((acc, span) => {
        acc[span.spanId] = acc[span.spanId] || {id: span.spanId, parentIds: [], data: span};
        acc[span.spanId].parentIds.push(span.parentSpanId);

        return acc;
      }, {});
  }, [testResultDetails]);

  const addSelected = useStoreActions(actions => actions.addSelectedElements);

  const handleOnSpanSelected = useCallback(
    (spanId: string) => {
      addSelected([{id: spanId}]);
      setSelectedSpan(spanMap[spanId]);
    },
    [addSelected, spanMap]
  );

  useGuidedTour(GuidedTours.Trace);

  useEffect(() => {
    let INTERVALID: any = null;

    INTERVALID = setInterval(() => {
      if (
        isError ||
        testResultDetails?.state === TestState.AWAITING_TRACE ||
        testResultDetails?.state === TestState.EXECUTING
      ) {
        refetchTrace();
      } else {
        INTERVALID && clearInterval(INTERVALID);
      }
    }, 1000);

    return () => INTERVALID && clearInterval(INTERVALID);
  }, [refetchTrace, testResultDetails?.state, isError]);

  useEffect(() => {
    if (testResultDetails && test && !isFirstLoad) {
      const resultList = runTest(testResultDetails.trace, test);

      setTraceResultList(resultList);

      updateTestResult({
        testId,
        resultId: testResultId,
        assertionResult: parseAssertionResultListToTestResult(resultList),
      });
    }
  }, [test]);

  useEffect(() => {
    if (testResultDetails && !isEmpty(testResultDetails.trace) && !testResultDetails?.assertionResult && test) {
      const resultList = runTest(testResultDetails.trace, test);

      setTraceResultList(resultList);
      setIsFirstLoad(false);

      updateTestResult({
        testId,
        resultId: testResultId,
        assertionResult: parseAssertionResultListToTestResult(resultList),
      });
    } else if (testResultDetails?.assertionResult && test) {
      setIsFirstLoad(false);
      setTraceResultList(
        parseTestResultToAssertionResultList(testResultDetails?.assertionResult, test, testResultDetails?.trace)
      );
    }
  }, [testResultDetails, test, testResultId, updateTestResult, testId]);

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
                  <TraceDiagram
                    spanMap={spanMap}
                    onSelectSpan={handleOnSpanSelected}
                    selectedSpan={selectedSpan}
                    trace={testResultDetails?.trace!}
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
                      <SpanDetail trace={testResultDetails?.trace} test={test} targetSpan={selectedSpan?.data} />
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
                        traceResultList={traceResultList}
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
