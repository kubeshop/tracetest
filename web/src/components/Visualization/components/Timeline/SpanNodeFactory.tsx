import {AxisScale} from '@visx/axis';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import {TNode} from 'types/Timeline.types';
import TestSpanNode from './TestSpanNode/TestSpanNode';
import TraceSpanNode from './TraceSpanNode/TraceSpanNode';

export interface IPropsComponent {
  index: number;
  indexParent: number;
  isCollapsed?: boolean;
  isMatched?: boolean;
  isSelected?: boolean;
  minStartTime: number;
  node: TNode;
  onClick(id: string): void;
  onCollapse(id: string): void;
  xScale: AxisScale;
}

const ComponentMap: Record<NodeTypesEnum, (props: IPropsComponent) => React.ReactElement> = {
  [NodeTypesEnum.TestSpan]: TestSpanNode,
  [NodeTypesEnum.TraceSpan]: TraceSpanNode,
};

interface IProps extends IPropsComponent {
  type: NodeTypesEnum;
}

const SpanNodeFactory = ({type, ...props}: IProps) => {
  const Component = ComponentMap[type];

  return <Component {...props} />;
};

export default SpanNodeFactory;
