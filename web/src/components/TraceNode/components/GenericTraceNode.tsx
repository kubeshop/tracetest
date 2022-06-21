import Text from 'antd/lib/typography/Text';
import {capitalize} from 'lodash';
import {Handle, Position} from 'react-flow-renderer';

// import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import * as S from '../TraceNode.styled';
import {TTraceNodeProps} from '../TraceNode';
import {SemanticGroupNames} from '../../../constants/SemanticGroupNames.constants';

const GenericTraceNode = ({id, data: span, selected}: TTraceNodeProps) => {
  // const {heading, primary} = SpanService.getSpanNodeInfo(span);
  // const spanTypeText = SemanticGroupNamesToText[type];

  return (
    <S.TraceNode selected={selected} data-cy={`trace-node-${span.label}`}>
      <S.TraceNotch spanType={SemanticGroupNames.Database}>
        <Text>{capitalize(span.label) || span.label}</Text>
      </S.TraceNotch>
      <Handle type="target" id={id} position={Position.Top} style={{top: 0, borderRadius: 0, visibility: 'hidden'}} />
      <S.TextContainer>
        {span.label && <S.NameText strong>{span.label}</S.NameText>}
        <S.NameText>{span.label}</S.NameText>
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

export default GenericTraceNode;
