import {NodeProps} from 'react-flow-renderer';

import useSpanData from 'hooks/useSpanData';
import {INodeDataSpan} from 'types/DAG.types';
import Footer from './Footer';
import SelectAsCurrent from './SelectAsCurrent';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';
import useSelectAsCurrent from '../../../hooks/useSelectAsCurrent';

interface IProps extends NodeProps<INodeDataSpan> {}

const TestSpanNode = ({data, id, selected}: IProps) => {
  const {span, testSpecs, testOutputs} = useSpanData(id);
  const {isLoading, onSelectAsCurrent, showSelectAsCurrent} = useSelectAsCurrent({
    selected,
    matched: data.isMatched,
    span,
  });

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
      {showSelectAsCurrent && <SelectAsCurrent isLoading={isLoading} onSelectAsCurrent={onSelectAsCurrent} />}
    </>
  );
};

export default TestSpanNode;
