import styled from 'styled-components';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
import Title from 'antd/lib/typography/Title';

import 'react-reflex/styles.css';

import {useState} from 'react';

import TraceDiagram from './TraceDiagram';
import TraceTimeline from './TraceTimeline';
import TraceData from './TraceData';

import data from './data.json';

const spanMap = data.data
  .map(i => i.spans)
  .flat()
  .reduce((acc: {[key: string]: {id: string; parentIds: string[]; data: any}}, span) => {
    acc[span.spanID] = acc[span.spanID] || {id: span.spanID, parentIds: [], data: span};
    span.references.forEach(p => {
      acc[span.spanID].parentIds.push(p.spanID);
    });
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
        <Title level={3}>{data.data[0].spans[0].operationName}</Title>
      </Header>
      <Grid>
        <ReflexContainer style={{minHeight: 1000, height: '100%'}} orientation="horizontal">
          <ReflexElement size={500} maxSize={500}>
            <ReflexContainer style={{height: '100%'}} orientation="vertical">
              <ReflexElement flex={0.5} className="left-pane">
                <div className="pane-content">
                  <TraceDiagram spanMap={spanMap} onSelectSpan={handleSelectSpan} selectedSpan={selectedSpan} />
                </div>
              </ReflexElement>
              <ReflexSplitter />
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
