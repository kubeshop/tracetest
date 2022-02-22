import styled from 'styled-components';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
import Title from 'antd/lib/typography/Title';

import 'react-reflex/styles.css';

import {useState} from 'react';

import TraceDiagram from './TraceDiagram';
import TraceTimeline from './TraceTimeline';
import TraceData from './TraceData';

import data from './data.json';

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

const Header = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  height: 64px;
  padding: 0 32px;
  border-bottom: 1px solid rgb(213, 215, 224);
`;

const Trace = () => {
  const [selectedSpan, setSelectedSpan] = useState<any>({});

  const handleSelectSpan = (span: any) => {
    setSelectedSpan(span);
  };

  return (
    <main>
      <Header>
        <Title level={3}>Title</Title>
      </Header>
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
                <div className="pane-content">
                  <TraceData json={JSON.parse(JSON.stringify(selectedSpan))} />
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
