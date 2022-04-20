import styled from 'styled-components';
import {useStoreActions} from 'react-flow-renderer';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
import {isEmpty} from 'lodash';

import {Button, Skeleton, Tabs} from 'antd';

import 'react-reflex/styles.css';

import {useCallback, useEffect, useMemo, useState} from 'react';
import {AssertionResult, ISpan} from 'types';
import {useGetTestByIdQuery, useGetTestResultByIdQuery, useUpdateTestResultMutation} from 'services/TestService';

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
  testId: string;
  testResultId: string;
};

const Trace: React.FC<TraceProps> = ({testId, testResultId}) => {
  const [selectedSpan, setSelectedSpan] = useState<TSpanInfo | undefined>();
  const [traceResultList, setTraceResultList] = useState<AssertionResult[]>([]);
  const [isFirstLoad, setIsFirstLoad] = useState(true);
  const [updateTestResult] = useUpdateTestResultMutation();
  const {data: test} = useGetTestByIdQuery(testId);

  const {
    data: testResultDetails,
    isLoading: isLoadingTrace,
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

  const handleReload = useCallback(() => {
    refetchTrace();
  }, [refetchTrace]);

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

  if (isLoadingTrace) {
    return <Skeleton />;
  }

  if (isError || Object.keys(testResultDetails?.trace || {}).length === 0) {
    return (
      <div>
        <Button onClick={handleReload} loading={isLoadingTrace}>
          Reload
        </Button>
      </div>
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
                        <SpanDetail trace={testResultDetails?.trace!} test={test} targetSpan={selectedSpan?.data!} />
                      </Tabs.TabPane>
                      <Tabs.TabPane tab="Test Results" key="2">
                        <TestResults
                          onSpanSelected={handleOnSpanSelected}
                          trace={testResultDetails?.trace!}
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
