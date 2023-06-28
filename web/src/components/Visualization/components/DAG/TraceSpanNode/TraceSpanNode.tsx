import {NodeProps} from 'react-flow-renderer';

import useSpanData from 'hooks/useSpanData';
import {INodeDataSpan} from 'types/DAG.types';
import AnalyzerErrors from './AnalyzerErrors';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';

interface IProps extends NodeProps<INodeDataSpan> {}

const TraceSpanNode = ({data, id, selected}: IProps) => {
  const {span, analyzerErrors} = useSpanData(id);

  return (
    <BaseSpanNode
      className={data.isMatched ? 'matched' : ''}
      footer={analyzerErrors && <AnalyzerErrors errors={analyzerErrors} />}
      id={id}
      isMatched={data.isMatched}
      isSelected={selected}
      span={span}
    />
  );
};

export default TraceSpanNode;
