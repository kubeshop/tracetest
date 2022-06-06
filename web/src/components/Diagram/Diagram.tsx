import {TestState} from '../../constants/TestRun.constants';
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
  affectedSpans: string[];
  onSelectSpan?(spanId: string): void;
  selectedSpan?: TSpan;
  trace: TTrace;
  type: SupportedDiagrams;
  runState: TTestRunState;
}

const ComponentMap: Record<string, typeof DAGComponent | typeof TimelineChart> = {
  [SupportedDiagrams.DAG]: DAGComponent,
  [SupportedDiagrams.Timeline]: TimelineChart,
};

const Diagram: React.FC<IDiagramProps> = ({type, runState, ...props}) => {
  const Component = ComponentMap[type || ''] || DAGComponent;

  return runState === TestState.FINISHED ? (
    <Component type={type} runState={runState} {...props} />
  ) : (
    <SkeletonDiagram onSelectSpan={props.onSelectSpan} />
  );
};

export default Diagram;
