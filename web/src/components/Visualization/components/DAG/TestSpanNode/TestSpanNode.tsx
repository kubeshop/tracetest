import {NodeProps} from 'react-flow-renderer';

import CurrentSpanSelector from 'components/CurrentSpanSelector';
import useSpanData from 'hooks/useSpanData';
import {INodeDataSpan} from 'types/DAG.types';
import Footer from './Footer';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';
import useShowSelectAsCurrent from '../../../hooks/useShowSelectAsCurrent';

interface IProps extends NodeProps<INodeDataSpan> {}

const TestSpanNode = ({data, id, selected}: IProps) => {
  const {span, testSpecs, testOutputs} = useSpanData(id);
  const showSelectAsCurrent = useShowSelectAsCurrent({selected, matched: data.isMatched});

  return (
    <>
      <BaseSpanNode
        className={`${data.isMatched && 'matched'} ${showSelectAsCurrent && 'selectedAsCurrent'}`}
        footer={<Footer testOutputs={testOutputs} testSpecs={testSpecs} />}
        id={id}
        isMatched={data.isMatched}
        isSelected={selected}
        span={span}
      />
      {showSelectAsCurrent && <CurrentSpanSelector spanId={id} />}
    </>
  );
};

export default TestSpanNode;
