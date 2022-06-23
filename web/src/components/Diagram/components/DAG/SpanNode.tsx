import Text from 'antd/lib/typography/Text';
import {capitalize} from 'lodash';
import {Handle, NodeProps, Position} from 'react-flow-renderer';

import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {INodeDataSpan} from 'types/DAG.types';
import * as S from './SpanNode.styled';

interface IProps extends NodeProps<INodeDataSpan> {}

const SpanNode = ({data, id, selected}: IProps) => {
  const spanTypeText = SemanticGroupNamesToText[data.type];
  const className = `${data.isAffected ? 'affected' : ''} ${data.isMatched ? 'matched' : ''}`;

  return (
    <S.Container className={className} data-cy={`trace-node-${data.type}`} selected={selected}>
      <S.Header type={data.type}>
        <Text>{capitalize(data.heading) || spanTypeText}</Text>
      </S.Header>
      <Handle id={id} position={Position.Top} style={{top: 0, visibility: 'hidden'}} type="target" />
      <S.Body>
        {data.primary && <S.BodyText strong>{data.primary}</S.BodyText>}
        <S.BodyText>{data.name}</S.BodyText>
      </S.Body>
      <Handle id={id} position={Position.Bottom} style={{bottom: 0, visibility: 'hidden'}} type="source" />
    </S.Container>
  );
};

export default SpanNode;
