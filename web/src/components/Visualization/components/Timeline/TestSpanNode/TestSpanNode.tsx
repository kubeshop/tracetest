import {NodeHeight} from 'constants/Timeline.constants';
import useSpanData from 'hooks/useSpanData';
import Header from './Header';
import SelectAsCurrent from './SelectAsCurrent';
import BaseSpanNode from '../BaseSpanNode/BaseSpanNode';
import {IPropsComponent} from '../SpanNodeFactory';
import useSelectAsCurrent from '../../../hooks/useSelectAsCurrent';

const TestSpanNode = (props: IPropsComponent) => {
  const {index, isMatched, isSelected, node} = props;
  const {span, testSpecs, testOutputs} = useSpanData(node.data.id);
  const {isLoading, onSelectAsCurrent, showSelectAsCurrent} = useSelectAsCurrent({
    selected: isSelected ?? false,
    matched: isMatched ?? false,
    span,
  });
  const positionTop = index * NodeHeight;

  return (
    <>
      <BaseSpanNode
        {...props}
        className={`${isMatched ? 'matched' : ''} ${showSelectAsCurrent ? 'selectedAsCurrent' : ''}`}
        header={
          <Header
            hasOutputs={!!testOutputs?.length}
            totalFailedChecks={testSpecs?.failed?.length}
            totalPassedChecks={testSpecs?.passed?.length}
          />
        }
        span={span}
      />
      {showSelectAsCurrent && (
        <SelectAsCurrent isLoading={isLoading} onSelectAsCurrent={onSelectAsCurrent} positionTop={positionTop} />
      )}
    </>
  );
};

export default TestSpanNode;
