import Text from 'antd/lib/typography/Text';
import {capitalize} from 'lodash';
import {Handle, NodeProps, Position} from 'react-flow-renderer';

import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {INodeDataSpan} from 'types/DAG.types';
import * as S from './TraceNode.styled';

interface IProps extends NodeProps<INodeDataSpan> {}

const TraceNode = ({data, id, selected}: IProps) => {
  const spanTypeText = SemanticGroupNamesToText[data.type];
  const className = `${data.isAffected ? 'affected' : ''} ${data.isMatched ? 'matched' : ''}`;

  return (
    <S.TraceNode className={className} data-cy={`trace-node-${data.type}`} selected={selected}>
      <S.TraceNotch spanType={data.type}>
        <Text>{capitalize(data.heading) || spanTypeText}</Text>
      </S.TraceNotch>
      <Handle id={id} position={Position.Top} style={{top: 0, visibility: 'hidden'}} type="target" />
      <S.TextContainer>
        {data.primary && <S.NameText strong>{data.primary}</S.NameText>}
        <S.NameText>{data.name}</S.NameText>
      </S.TextContainer>
      <Handle id={id} position={Position.Bottom} style={{bottom: 0, visibility: 'hidden'}} type="source" />
    </S.TraceNode>
  );
};

export default TraceNode;
