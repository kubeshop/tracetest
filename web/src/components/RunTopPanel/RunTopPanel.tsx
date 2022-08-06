import Diagram, {SupportedDiagrams} from 'components/Diagram/Diagram';
import {useSpan} from 'providers/Span/Span.provider';
import {useState} from 'react';
import {ReflexSplitter} from 'react-reflex';
import {TTestRun} from 'types/TestRun.types';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
import DiagramSwitcher from '../DiagramSwitcher';
import {useRunLayout} from '../RunLayout';
import SpanDetail from '../SpanDetail';
import * as S from './RunTopPanel.styled';

interface IProps {
  run: TTestRun;
}

const RunTopPanel = ({run}: IProps) => {
  const [diagramType, setDiagramType] = useState<SupportedDiagrams>(SupportedDiagrams.DAG);
  const {onSearch, selectedSpan} = useSpan();
  const {isBottomPanelOpen} = useRunLayout();

  return (
    <S.Container orientation="vertical">
      <S.LeftPanel>
        <div style={{height: 'calc( 100% - 48px )', padding: 24}}>
          <DiagramSwitcher
            onTypeChange={type => {
              TraceAnalyticsService.onSwitchDiagramView(type);
              setDiagramType(type);
            }}
            onSearch={onSearch}
            selectedType={diagramType}
          />
          <Diagram trace={run.trace!} runState={run.state} type={diagramType} />
        </div>
      </S.LeftPanel>
      {isBottomPanelOpen ? null : <ReflexSplitter />}
      <S.RightPanel>
        <SpanDetail span={selectedSpan} />
      </S.RightPanel>
    </S.Container>
  );
};

export default RunTopPanel;
