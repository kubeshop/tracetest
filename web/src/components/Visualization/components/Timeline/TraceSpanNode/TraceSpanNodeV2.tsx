import useSpanData from 'hooks/useSpanData';
// import Header from './Header';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNodeV2';
import {IPropsComponent} from '../SpanNodeFactoryV2';

const TraceSpanNode = (props: IPropsComponent) => {
  const {node} = props;
  const {span} = useSpanData(node.data.id);

  return (
    <BaseSpanNode
      {...props}
      // header={<Header hasAnalyzerErrors={!!analyzerErrors} />}
      span={span}
    />
  );
};

export default TraceSpanNode;
