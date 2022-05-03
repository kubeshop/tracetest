import Text from 'antd/lib/typography/Text';
import {upperCase} from 'lodash';
import React from 'react';
import {Handle, NodeProps, Position} from 'react-flow-renderer';
import {SemanticGroupNamesToText} from '../../constants/SemanticGroupNames.constants';
import {TSpan} from '../../types/Span.types';
import * as S from './TraceDiagram.styled';
import SpanService from '../../services/Span.service';
import { TTrace } from '../../types/Trace.types';

type TTraceNodeProps = NodeProps<{span: TSpan; trace: TTrace}>;

const TraceNode: React.FC<TTraceNodeProps> = ({
  id,
  data: {
    span: {name, type},
    span,
  },
  selected,
}) => {
  const {heading, primary} = SpanService.getSpanNodeInfo(span);
  const spanTypeText = SemanticGroupNamesToText[type];

  return (
    <S.TraceNode selected={selected}>
      <S.TraceNotch spanType={type}>
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
