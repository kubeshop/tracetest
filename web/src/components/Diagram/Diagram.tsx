import SkeletonDiagram from 'components/SkeletonDiagram';
import {TestState} from 'constants/TestRun.constants';
import {DAGProvider} from 'providers/DAG';
import {useSpan} from 'providers/Span/Span.provider';
import {TSpan} from 'types/Span.types';
import {TTestRunState} from 'types/TestRun.types';
import {TTrace} from 'types/Trace.types';
import DAGComponent from './components/DAG';
import TimelineComponent from './components/Timeline';

export enum SupportedDiagrams {
  DAG = 'dag',
  Timeline = 'timeline',
}

export interface IProps {
  runState: TTestRunState;
  trace: TTrace;
  type: SupportedDiagrams;
}

export interface IDiagramComponentProps {
  affectedSpans: string[];
  matchedSpans: string[];
  onSelectSpan(spanId: string): void;
  selectedSpan?: TSpan;
  spanList: TSpan[];
}

const ComponentMap: Record<string, typeof DAGComponent | typeof TimelineComponent> = {
  [SupportedDiagrams.DAG]: DAGComponent,
  [SupportedDiagrams.Timeline]: TimelineComponent,
};

const Diagram = ({runState, trace, type}: IProps) => {
  const Component = ComponentMap[type || ''] || DAGComponent;

  const {onClearAffectedSpans, onClearSelectedSpan, onSelectSpan, selectedSpan, affectedSpans, matchedSpans} =
    useSpan();
  const spanList = trace.spans || [];

  return runState === TestState.FINISHED ? (
    <DAGProvider spans={spanList}>
      <Component {...{spanList, onSelectSpan, selectedSpan, affectedSpans, matchedSpans}} />
    </DAGProvider>
  ) : (
    <SkeletonDiagram onClearAffectedSpans={onClearAffectedSpans} onClearSelectedSpan={onClearSelectedSpan} />
  );
};

export default Diagram;
