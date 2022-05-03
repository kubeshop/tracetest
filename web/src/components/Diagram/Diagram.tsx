import { ISpan } from '../../types/Span.types';
import { ITrace } from '../../types/Trace.types';
import DAGComponent from './components/DAG';

export enum SupportedDiagrams {
  DAG = 'dag',
}

export interface IDiagramProps {
  type: SupportedDiagrams;
  trace: ITrace;
  selectedSpan?: ISpan;
  onSelectSpan?(spanId: string): void;
}

const ComponentMap: Record<string, typeof DAGComponent> = {
  [SupportedDiagrams.DAG]: DAGComponent,
};

const Diagram: React.FC<IDiagramProps> = ({type, ...props}) => {
  const Component = ComponentMap[type || ''] || DAGComponent;

  return <Component type={type} {...props} />;
};

export default Diagram;
