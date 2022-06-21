import {useState} from 'react';

import Diagram from 'components/Diagram';
import {SupportedDiagrams} from 'components/Diagram/Diagram';
import DiagramSwitcher from 'components/DiagramSwitcher';
import SpanDetail from 'components/SpanDetail';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import {useSpan} from 'providers/Span/Span.provider';
import {TTestRun} from 'types/TestRun.types';
import * as S from './RunTopPanel.styled';

interface IProps {
  run: TTestRun;
}

const RunTopPanel = ({run}: IProps) => {
  const [diagramType, setDiagramType] = useState<SupportedDiagrams>(SupportedDiagrams.DAG);
  const {onSearch, selectedSpan} = useSpan();

  return (
    <S.Container>
      <S.LeftPanel>
        <DiagramSwitcher
          onTypeChange={type => {
            TraceAnalyticsService.onSwitchDiagramView(type);
            setDiagramType(type);
          }}
          onSearch={onSearch}
          selectedType={diagramType}
        />
        <Diagram trace={run.trace!} runState={run.state} type={diagramType} />
      </S.LeftPanel>
      <S.RightPanel>
        <SpanDetail span={selectedSpan} />
      </S.RightPanel>
    </S.Container>
  );
};

export default RunTopPanel;
