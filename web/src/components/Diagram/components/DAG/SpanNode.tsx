import {ClockCircleOutlined} from '@ant-design/icons';
import {Handle, NodeProps, Position} from 'react-flow-renderer';

import {SemanticGroupNamesToIcon} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {INodeDataSpan} from 'types/DAG.types';
import * as S from './SpanNode.styled';

interface IProps extends NodeProps<INodeDataSpan> {}

const SpanNode = ({data, id, selected}: IProps) => {
  const spansResult = useAppSelector(TestDefinitionSelectors.selectSpansResult);
  const className = `${data.isAffected ? 'affected' : ''} ${data.isMatched ? 'matched' : ''}`;
  const Icon = SemanticGroupNamesToIcon[data.type];

  return (
    <S.Container
      className={className}
      data-cy={`trace-node-${data.type}`}
      $affected={data.isAffected}
      $selected={selected}
    >
      <Handle id={id} position={Position.Top} style={{top: 0, visibility: 'hidden'}} type="target" />

      <S.Header $type={data.type}>
        <S.HeaderText>{data.name}</S.HeaderText>
        <S.IconContainer $type={data.type}>
          <Icon />
        </S.IconContainer>
      </S.Header>

      <S.Body>
        <S.BodyText>
          {data.serviceName} {SpanKindToText[data.kind]}
        </S.BodyText>
        <S.BodyText $secondary>{data.system}</S.BodyText>
        <S.Badge>
          <ClockCircleOutlined />
          <S.BadgeText>{data.duration}ms</S.BadgeText>
        </S.Badge>
      </S.Body>

      <S.Footer>
        {Boolean(spansResult[data.id]?.passed) && <S.Dot $type="success" />}
        {Boolean(spansResult[data.id]?.failed) && <S.Dot $type="error" />}
      </S.Footer>

      <Handle id={id} position={Position.Bottom} style={{bottom: 0, visibility: 'hidden'}} type="source" />
    </S.Container>
  );
};

export default SpanNode;
