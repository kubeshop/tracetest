import useSpanData from 'hooks/useSpanData';
import Header from './Header';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';
import {IPropsComponent} from '../SpanNodeFactory';

const TraceSpanNode = (props: IPropsComponent) => {
  const {node} = props;
  const {span, analyzerErrors} = useSpanData(node.data.id);

  return <BaseSpanNode {...props} header={<Header hasAnalyzerErrors={!!analyzerErrors} />} span={span} />;
};

export default TraceSpanNode;
