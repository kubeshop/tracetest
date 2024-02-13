import useSpanData from 'hooks/useSpanData';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';
import {IPropsComponent} from '../SpanNodeFactory';

const TraceSpanNode = (props: IPropsComponent) => {
  const {node} = props;
  const {span} = useSpanData(node.data.id);

  return <BaseSpanNode {...props} span={span} />;
};

export default TraceSpanNode;
