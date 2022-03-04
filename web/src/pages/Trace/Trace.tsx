import styled from 'styled-components';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';

import {Tabs} from 'antd';
import Text from 'antd/lib/typography/Text';

import 'react-reflex/styles.css';

import {useState} from 'react';

import TraceDiagram from './TraceDiagram';
import TraceTimeline from './TraceTimeline';
import TraceData from './TraceData';

import data from './data.json';
import AssertionList from './AssertionsList';
import {Test} from '../../types';

const spanMap = data.resourceSpans
  .map((i: any) => i.instrumentationLibrarySpans.map((el: any) => el.spans))
  .flat(2)
  .reduce((acc: {[key: string]: {id: string; parentIds: string[]; data: any}}, span) => {
    acc[span.spanId] = acc[span.spanId] || {id: span.spanId, parentIds: [], data: span};
    acc[span.spanId].parentIds.push(span.parentSpanId);

    return acc;
  }, {});

const Grid = styled.div`
  display: grid;
`;

const Trace = ({test}: {test: Test}) => {
  const [selectedSpan, setSelectedSpan] = useState<any>({});

  const handleSelectSpan = (span: any) => {
    setSelectedSpan(span);
  };

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
                <div className="pane-content" style={{paddingLeft: 8}}>
                  <div>
                    <Text>Service</Text>
                  </div>
                  <Tabs>
                    <Tabs.TabPane tab="Raw Data" key="1">
                      <TraceData json={JSON.parse(JSON.stringify(selectedSpan))} />
                    </Tabs.TabPane>
                    {spanMap[selectedSpan.id]?.data && (
                      <Tabs.TabPane tab="Assertions" key="2">
                        <AssertionList testId={test.id} targetSpan={spanMap[selectedSpan.id]?.data} />
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
              <TraceTimeline onSelectSpan={handleSelectSpan} selectedSpan={selectedSpan} />
            </div>
          </ReflexElement>
        </ReflexContainer>
      </Grid>
    </main>
  );
};

export default Trace;
