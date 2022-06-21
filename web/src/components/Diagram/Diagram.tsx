import {TestState} from '../../constants/TestRun.constants';
import {useSpan} from '../../providers/Span/Span.provider';
import {TSpan} from '../../types/Span.types';
import {TTestRunState} from '../../types/TestRun.types';
import {TTrace} from '../../types/Trace.types';
import SkeletonDiagram from '../SkeletonDiagram';
import DAGComponent from './components/DAG';
import {TimelineChart} from './components/TimelineChart';

export enum SupportedDiagrams {
  DAG = 'dag',
  Timeline = 'timeline',
}

export interface IDiagramProps {
  trace: TTrace;
  type: SupportedDiagrams;
  runState: TTestRunState;
}

export interface IDiagramComponentProps {
  spanList: TSpan[];
  affectedSpans: string[];
  matchedSpans: string[];
  selectedSpan?: TSpan;
  onSelectSpan(spanId: string): void;
}

const ComponentMap: Record<string, typeof DAGComponent | typeof TimelineChart> = {
  [SupportedDiagrams.DAG]: DAGComponent,
  [SupportedDiagrams.Timeline]: TimelineChart,
};

const Diagram: React.FC<IDiagramProps> = ({type, runState, trace}) => {
  const Component = ComponentMap[type || ''] || DAGComponent;
  const {onSelectSpan, selectedSpan, affectedSpans, matchedSpans} = useSpan();
  const spanList = trace.spans || [];

  return runState === TestState.FINISHED ? (
    <Component {...{spanList, onSelectSpan, selectedSpan, affectedSpans, matchedSpans}} />
  ) : (
    <SkeletonDiagram onSelectSpan={onSelectSpan} selectedSpan={selectedSpan} />
  );
};

export default Diagram;
