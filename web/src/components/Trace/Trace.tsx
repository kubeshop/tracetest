import {useCallback, useEffect, useMemo, useState} from 'react';
import styled from 'styled-components';
import {useNavigate} from 'react-router-dom';

import {useStoreActions} from 'react-flow-renderer';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
import {isEmpty} from 'lodash';

import {Button, Tabs, Typography} from 'antd';
import {CloseCircleFilled} from '@ant-design/icons';

import 'react-reflex/styles.css';

import {AssertionResult, ISpan, TestState} from 'types';

import {
  useGetTestByIdQuery,
  useGetTestResultByIdQuery,
  useRunTestMutation,
  useUpdateTestResultMutation,
} from 'services/TestService';
import {
  parseAssertionResultListToTestResult,
  parseTestResultToAssertionResultList,
  runTest,
} from 'services/TraceService';
import TraceDiagram from './TraceDiagram';
import TraceTimeline from './TraceTimeline';
import * as S from './Trace.styled';

import SpanDetail from './SpanDetail';
import TestResults from './TestResults';

const Grid = styled.div`
  display: grid;
`;

export type TSpanInfo = {
  id: string;
  parentIds: string[];
  data: ISpan;
};

type TSpanMap = Record<string, TSpanInfo>;

type TraceProps = {
  testId: string;
  testResultId: string;
  onDismissTrace: () => void;
};

const Trace: React.FC<TraceProps> = ({testId, testResultId, onDismissTrace}) => {
  const navigate = useNavigate();
  const [selectedSpan, setSelectedSpan] = useState<TSpanInfo | undefined>();
  const [traceResultList, setTraceResultList] = useState<AssertionResult[]>([]);
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

  useEffect(() => {
    let TIMEOUTID: any = null;
    if (
      isError ||
      testResultDetails?.state === TestState.AWAITING_TRACE ||
      testResultDetails?.state === TestState.EXECUTING
    ) {
      TIMEOUTID = setTimeout(() => {
        refetchTrace();
      }, 500);
    }
    return () => TIMEOUTID && clearTimeout(TIMEOUTID);
  }, [refetchTrace, testResultDetails, isError]);

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
    navigate(`/test/${result.testId}?resultId=${result.resultId}`);
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
              <ReflexElement flex={0.5} className="left-pane">
                <div className="pane-content">
                  <TraceDiagram spanMap={spanMap} onSelectSpan={handleOnSpanSelected} selectedSpan={selectedSpan} />
                </div>
              </ReflexElement>
              <ReflexElement flex={0.5} className="right-pane">
                <div className="pane-content" style={{padding: '14px 24px', overflow: 'hidden'}}>
                  <S.TraceTabs>
                    <Tabs.TabPane tab="Span detail" key="1">
                      <SpanDetail trace={testResultDetails?.trace} test={test} targetSpan={selectedSpan?.data} />
                    </Tabs.TabPane>
                    <Tabs.TabPane tab="Test Results" key="2">
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
                trace={testResultDetails?.trace}
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
