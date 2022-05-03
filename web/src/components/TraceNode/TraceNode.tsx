import {NodeProps} from 'react-flow-renderer';
import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';
import {TSpan} from '../../types/Span.types';
import GenericTraceNode from './components/GenericTraceNode';

export type TTraceNodeProps = NodeProps<TSpan>;

const ComponentMap: Record<string, typeof GenericTraceNode> = {
  [SemanticGroupNames.Http]: GenericTraceNode,
};

const TraceNode: React.FC<TTraceNodeProps> = ({data: span, ...props}) => {
  const Component = ComponentMap[span?.type || ''] || GenericTraceNode;

  return <Component data={span} {...props} />;
};

export default TraceNode;
