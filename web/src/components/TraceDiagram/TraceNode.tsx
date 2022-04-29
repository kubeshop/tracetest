import Text from 'antd/lib/typography/Text';
import {upperCase} from 'lodash';
import React from 'react';
import {Handle, NodeProps, Position} from 'react-flow-renderer';
import {SemanticGroupNames, SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import {getSpanNodeInfo} from '../../services/Span.service';
import {ISpan} from '../../types/Span.types';
import {ITrace} from '../../types/Trace.types';
import * as S from './TraceDiagram.styled';

type TTraceNodeProps = NodeProps<{span: ISpan; trace: ITrace}>;

const TraceNode: React.FC<TTraceNodeProps> = ({
  id,
  data: {
    span: {name, spanId},
    trace,
  },
  selected,
}) => {
  const {heading, spanType = SemanticGroupNames.General, primary} = getSpanNodeInfo(spanId, trace);
  const spanTypeText = SemanticGroupNamesToText[spanType];

  return (
    <S.TraceNode selected={selected}>
      <S.TraceNotch spanType={spanType}>
        <Text>{upperCase(heading || spanTypeText)}</Text>
      </S.TraceNotch>
      <Handle type="target" id={id} position={Position.Top} style={{top: 0, borderRadius: 0, visibility: 'hidden'}} />
      <S.TextContainer>
        {primary && <S.NameText strong>{primary}</S.NameText>}
        <S.NameText>{name}</S.NameText>
      </S.TextContainer>
      <Handle
        type="source"
        position={Position.Bottom}
        id={id}
        style={{bottom: 0, borderRadius: 0, visibility: 'hidden'}}
      />
    </S.TraceNode>
  );
};

export default TraceNode;
