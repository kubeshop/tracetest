import styled from 'styled-components';
import {ReflexContainer, ReflexSplitter, ReflexElement} from 'react-reflex';
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

const Trace = () => {
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
