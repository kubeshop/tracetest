import {ISpan} from '../../types/Span.types';
import {ITrace} from '../../types/Trace.types';
import DAGComponent from './components/DAG';
import {TimelineChart} from './components/TimelineChart';

export enum SupportedDiagrams {
  DAG = 'dag',
  Timeline = 'timeline',
}

export interface IDiagramProps {
  type: SupportedDiagrams;
  trace: ITrace;
  selectedSpan?: ISpan;
  onSelectSpan?(spanId: string): void;
}

const ComponentMap: Record<string, typeof DAGComponent | typeof TimelineChart> = {
  [SupportedDiagrams.DAG]: DAGComponent,
  [SupportedDiagrams.Timeline]: TimelineChart,
};

const Diagram: React.FC<IDiagramProps> = ({type, ...props}) => {
  const Component = ComponentMap[type || ''] || DAGComponent;

  return <Component type={type} {...props} />;
};

export default Diagram;
