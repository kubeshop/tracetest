import useSpanData from 'hooks/useSpanData';
import Header from './Header';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';
import {IPropsComponent} from '../SpanNodeFactory';

const TestSpanNode = (props: IPropsComponent) => {
  const {node} = props;
  const {span, testSpecs} = useSpanData(node.data.id);

  return (
    <BaseSpanNode
      {...props}
      header={<Header totalFailedChecks={testSpecs?.failed?.length} totalPassedChecks={testSpecs?.passed?.length} />}
      span={span}
    />
  );
};

export default TestSpanNode;
