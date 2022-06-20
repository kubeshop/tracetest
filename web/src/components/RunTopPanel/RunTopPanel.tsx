import {useState} from 'react';

import Diagram from 'components/Diagram';
import {SupportedDiagrams} from 'components/Diagram/Diagram';
import DiagramSwitcher from 'components/DiagramSwitcher';
import SpanDetail from 'components/SpanDetail';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import {useSpan} from 'providers/Span/Span.provider';
import {TSpan} from 'types/Span.types';
import {TTestRun} from 'types/TestRun.types';
import * as S from './RunTopPanel.styled';

interface IProps {
  onSelectSpan: (spanId: string) => void;
  run: TTestRun;
  selectedSpan: TSpan;
}

const RunTopPanel = ({onSelectSpan, run, selectedSpan}: IProps) => {
  const [diagramType, setDiagramType] = useState<SupportedDiagrams>(SupportedDiagrams.DAG);
  const {affectedSpans} = useSpan();

  return (
    <S.Container>
      <S.LeftPanel>
        <DiagramSwitcher
          onTypeChange={type => {
            TraceAnalyticsService.onSwitchDiagramView(type);
            setDiagramType(type);
          }}
          onSearch={() => console.log('onSearch')}
          selectedType={diagramType}
        />
        <Diagram
          affectedSpans={affectedSpans}
          onSelectSpan={onSelectSpan}
          selectedSpan={selectedSpan}
          trace={run.trace!}
          runState={run.state}
          type={diagramType}
        />
      </S.LeftPanel>
      <S.RightPanel>
        <SpanDetail span={selectedSpan} />
      </S.RightPanel>
    </S.Container>
  );
};

export default RunTopPanel;
