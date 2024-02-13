import {NodeTypesEnum} from 'constants/Visualization.constants';
import {TNode} from 'types/Timeline.types';
// import TestSpanNode from './TestSpanNode/TestSpanNode';
import TraceSpanNode from './TraceSpanNode/TraceSpanNode';

export interface IPropsComponent {
  index: number;
  node: TNode;
  style: React.CSSProperties;
}

const ComponentMap: Record<NodeTypesEnum, (props: IPropsComponent) => React.ReactElement> = {
  [NodeTypesEnum.TestSpan]: TraceSpanNode,
  [NodeTypesEnum.TraceSpan]: TraceSpanNode,
};

interface IProps {
  data: TNode[];
  index: number;
  style: React.CSSProperties;
}

const SpanNodeFactory = ({data, ...props}: IProps) => {
  const node = data[props.index];
  const Component = ComponentMap[node.type];

  return <Component {...props} node={node} />;
};

export default SpanNodeFactory;
