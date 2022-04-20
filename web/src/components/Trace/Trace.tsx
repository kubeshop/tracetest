import {useCallback, useEffect, useMemo, useState} from 'react';
import styled from 'styled-components';
import {useStoreActions} from 'react-flow-renderer';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';

import {Button, Skeleton, Tabs, Typography} from 'antd';
import {CloseCircleFilled} from '@ant-design/icons';

import 'react-reflex/styles.css';

import {AssertionResult, ISpan, Test, TestState} from 'types';
import {useGetTestResultByIdQuery, useUpdateTestResultMutation} from 'services/TestService';

import TraceDiagram from './TraceDiagram';
import TraceTimeline from './TraceTimeline';
import * as S from './Trace.styled';

import SpanDetail from './SpanDetail';
import TestResults from './TestResults';
import {
  parseAssertionResultListToTestResult,
  parseTestResultToAssertionResultList,
  runTest,
} from '../../services/TraceService';

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
  test: Test;
  testResultId: string;
};

const Trace: React.FC<TraceProps> = ({test, testResultId}) => {
  const [selectedSpan, setSelectedSpan] = useState<TSpanInfo | undefined>();
  const [traceResultList, setTraceResultList] = useState<AssertionResult[]>([]);
  const [updateTestResult] = useUpdateTestResultMutation();

  const {
    data: testResultDetails,
    isLoading: isLoadingTrace,
    isError,
    refetch: refetchTrace,
  } = useGetTestResultByIdQuery({testId: test.testId, resultId: testResultId});

  const spanMap = useMemo(() => {
    const spanList: ISpan[] =
      testResultDetails?.trace?.resourceSpans
        ?.map(i => i.instrumentationLibrarySpans.map((el: any) => el.spans))
        ?.flat(2) || [];

    const defaultSpan = spanList.find(({parentSpanId}) => !parentSpanId);

    return spanList.reduce<TSpanMap>((acc, span) => {
      if (span.spanId === defaultSpan?.spanId) return acc;

      acc[span.spanId] = acc[span.spanId] || {id: span.spanId, parentIds: [], data: span};
      if (span.parentSpanId !== defaultSpan?.spanId && span.parentSpanId)
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

  const handleReload = useCallback(() => {
    refetchTrace();
  }, [refetchTrace]);

  useEffect(() => {
    if (testResultDetails && !testResultDetails?.assertionResult) {
      const resultList = runTest(testResultDetails.trace, test);

      setTraceResultList(resultList);

      updateTestResult({
        testId: test.testId,
        resultId: testResultId,
        assertionResult: parseAssertionResultListToTestResult(resultList),
      });
    } else if (testResultDetails?.assertionResult) {
      setTraceResultList(
        parseTestResultToAssertionResultList(testResultDetails?.assertionResult, test, testResultDetails?.trace)
      );
    }
  }, [testResultDetails, test, testResultId, updateTestResult]);

  useEffect(() => {
    if (testResultDetails) {
      const resultList = runTest(testResultDetails.trace, test);

      setTraceResultList(resultList);

      updateTestResult({
        testId: test.testId,
        resultId: testResultId,
        assertionResult: parseAssertionResultListToTestResult(resultList),
      });
    }
  }, [test, testResultDetails, testResultId, updateTestResult]);

  if (isLoadingTrace) {
    return <Skeleton />;
  }

  if (
    isError ||
    Object.keys(testResultDetails?.trace || {}).length === 0 ||
    testResultDetails?.state === TestState.FAILED
  ) {
    return (
      <S.FailedTrace>
        <CloseCircleFilled style={{color: 'red', fontSize: 32}} />
        <Typography.Title level={2}>Test Run Failed</Typography.Title>
        <div style={{display: 'grid', gap: 8, gridTemplateColumns: '1fr 1fr'}}>
          <Button>Rerun Test</Button>
          <Button>Cancel</Button>
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
                  {Boolean(selectedSpan) && (
                    <S.TraceTabs>
                      <Tabs.TabPane tab="Span detail" key="1">
                        <SpanDetail
                          trace={testResultDetails?.trace!}
                          testId={test.testId}
                          targetSpan={selectedSpan?.data!}
                        />
                      </Tabs.TabPane>
                      <Tabs.TabPane tab="Test Results" key="2">
                        <TestResults
                          onSpanSelected={handleOnSpanSelected}
                          trace={testResultDetails?.trace!}
                          test={test}
                          traceResultList={traceResultList}
                        />
                      </Tabs.TabPane>
                    </S.TraceTabs>
                  )}
                </div>
              </ReflexElement>
            </ReflexContainer>
          </ReflexElement>
          <ReflexSplitter />
          <ReflexElement>
            <div className="pane-content">
              {testResultDetails && selectedSpan && (
                <TraceTimeline
                  trace={testResultDetails?.trace}
                  onSelectSpan={handleOnSpanSelected}
                  selectedSpan={selectedSpan}
                />
              )}
            </div>
          </ReflexElement>
        </ReflexContainer>
      </Grid>
    </main>
  );
};

export default Trace;
