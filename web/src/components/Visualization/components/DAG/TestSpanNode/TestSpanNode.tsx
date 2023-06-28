import {NodeProps} from 'react-flow-renderer';

import CurrentSpanSelector from 'components/CurrentSpanSelector';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import useSpanData from 'hooks/useSpanData';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {INodeDataSpan} from 'types/DAG.types';
import Footer from './Footer';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';

interface IProps extends NodeProps<INodeDataSpan> {}

const TestSpanNode = ({data, id, selected}: IProps) => {
  const {span, testSpecs, testOutputs} = useSpanData(id);

  const {matchedSpans} = useSpan();
  const {isOpen: isTestSpecFormOpen} = useTestSpecForm();
  const {isOpen: isTestOutputFormOpen} = useTestOutput();

  const showSelectAsCurrent =
    selected && !data.isMatched && !!matchedSpans.length && (isTestSpecFormOpen || isTestOutputFormOpen);

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
