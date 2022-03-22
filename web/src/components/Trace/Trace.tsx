import styled from 'styled-components';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';

import {Button, Skeleton, Tabs} from 'antd';
import Text from 'antd/lib/typography/Text';

import 'react-reflex/styles.css';

import {useMemo, useState} from 'react';
import {Test} from 'types';
import {useGetTestResultByIdQuery} from 'services/TestService';

import TraceDiagram from './TraceDiagram';
import TraceTimeline from './TraceTimeline';
import TraceData from './TraceData';

import AssertionList from './AssertionsList';

const Grid = styled.div`
  display: grid;
`;

const Trace = ({test, testResultId}: {test: Test; testResultId: string}) => {
  const [selectedSpan, setSelectedSpan] = useState<any>({});
  const {
    data: traceData,
    isLoading: isLoadingTrace,
    isError,
    refetch: refetchTrace,
  } = useGetTestResultByIdQuery({testId: test.testId, resultId: testResultId});

  const spanMap = useMemo(() => {
    return traceData?.trace?.resourceSpans
      ?.map(i => i.instrumentationLibrarySpans.map((el: any) => el.spans))
      ?.flat(2)
      ?.reduce((acc, span) => {
        acc[span.spanId] = acc[span.spanId] || {id: span.spanId, parentIds: [], data: span};
        acc[span.spanId].parentIds.push(span.parentSpanId);

        return acc;
      }, {});
  }, [traceData]);

  const handleSelectSpan = (span: any) => {
    setSelectedSpan(span);
  };

  const handleReload = () => {
    refetchTrace();
  };

  if (isLoadingTrace) {
    return <Skeleton />;
  }

  if (isError) {
    return (
      <div>
        <Button onClick={handleReload}>Reload</Button>
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
                  <TraceDiagram spanMap={spanMap} onSelectSpan={handleSelectSpan} selectedSpan={selectedSpan} />
                </div>
              </ReflexElement>

              <ReflexElement flex={0.5} className="right-pane">
                <div className="pane-content" style={{paddingLeft: 8, overflow: 'hidden'}}>
                  <div>
                    <Text>Service</Text>
                  </div>
                  <Tabs>
                    <Tabs.TabPane tab="Raw Data" key="1">
                      <TraceData json={JSON.parse(JSON.stringify(selectedSpan))} />
                    </Tabs.TabPane>
                    {spanMap[selectedSpan.id]?.data && (
                      <Tabs.TabPane tab="Assertions" key="2">
                        <AssertionList
                          trace={traceData?.trace}
                          testId={test.testId}
                          targetSpan={spanMap[selectedSpan.id]?.data}
                        />
                      </Tabs.TabPane>
                    )}
                  </Tabs>
                </div>
              </ReflexElement>
            </ReflexContainer>
          </ReflexElement>
          <ReflexSplitter />
          <ReflexElement>
            <div className="pane-content">
              {traceData && (
                <TraceTimeline trace={traceData?.trace} onSelectSpan={handleSelectSpan} selectedSpan={selectedSpan} />
              )}
            </div>
          </ReflexElement>
        </ReflexContainer>
      </Grid>
    </main>
  );
};

export default Trace;
