import PanelLayout from 'components/ResizablePanels';
import {useMemo} from 'react';
import TestRun from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunDetailTrace.styled';
import SetupAlert from '../SetupAlert';
import {getAnalyzerPanel} from './AnalyzerPanel';
import {getSpanDetailsPanel} from './SpanDetailsPanel';
import {geTracePanel} from './TracePanel';

interface IProps {
  run: TestRun;
  runEvents: TestRunEvent[];
  testId: string;
}

export enum VisualizationType {
  Dag,
  Timeline,
}

const RunDetailTrace = ({run, runEvents, testId}: IProps) => {
  const panels = useMemo(
    () => [getSpanDetailsPanel(testId, run), geTracePanel(testId, run, runEvents), getAnalyzerPanel(run)],
    [run, runEvents, testId]
  );

  return (
    <S.Container>
      <SetupAlert />
      <PanelLayout panels={panels} />
    </S.Container>
  );
};

export default RunDetailTrace;
