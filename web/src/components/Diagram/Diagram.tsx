import {TestState} from 'constants/TestRun.constants';
import {DAGProvider} from 'providers/DAG';
import {useSpan} from 'providers/Span/Span.provider';
import {TSpan} from 'types/Span.types';
import {TTestRunState} from 'types/TestRun.types';
import {TTrace} from 'types/Trace.types';
import DAGComponent from './components/DAG';
import {TimelineChart} from './components/TimelineChart';
import SkeletonDiagram from '../SkeletonDiagram';

export enum SupportedDiagrams {
  DAG = 'dag',
  Timeline = 'timeline',
}

export interface IProps {
  trace: TTrace;
  type: SupportedDiagrams;
  runState: TTestRunState;
}

export interface IDiagramComponentProps {
  affectedSpans: string[];
  matchedSpans: string[];
  onSelectSpan(spanId: string): void;
  selectedSpan?: TSpan;
  spanList: TSpan[];
}

const ComponentMap: Record<string, typeof DAGComponent | typeof TimelineChart> = {
  [SupportedDiagrams.DAG]: DAGComponent,
  [SupportedDiagrams.Timeline]: TimelineChart,
};

const Diagram = ({type, runState, trace}: IProps) => {
  const Component = ComponentMap[type || ''] || DAGComponent;
  const {onClearAffectedSpans, onClearSelectedSpan, onSelectSpan, selectedSpan, affectedSpans, matchedSpans} =
    useSpan();
  const spanList = trace.spans || [];

  return runState === TestState.FINISHED ? (
    <DAGProvider>
      <Component {...{spanList, onSelectSpan, selectedSpan, affectedSpans, matchedSpans}} />
    </DAGProvider>
  ) : (
    <SkeletonDiagram onClearAffectedSpans={onClearAffectedSpans} onClearSelectedSpan={onClearSelectedSpan} />
  );
};

export default Diagram;
