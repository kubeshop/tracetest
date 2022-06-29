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
    <S.Container className={className} data-cy={`trace-node-${data.type}`} $selected={selected} $type={data.type}>
      <Handle id={id} position={Position.Top} style={{top: 0, visibility: 'hidden'}} type="target" />
      <S.Header>
        <S.Logo $type={data.type}>
          <Icon />
        </S.Logo>
        <S.HeaderText ellipsis={{rows: 2}} strong $type={data.type}>
          {data.name}
        </S.HeaderText>
      </S.Header>
      <S.Body>
        <S.BodyText>
          {data.serviceName} {SpanKindToText[data.kind]}
        </S.BodyText>

        <S.Badge>
          <ClockCircleOutlined style={{color: '#031849', marginRight: 4}} />
          <S.BadgeText>{data.duration}ms</S.BadgeText>
        </S.Badge>
      </S.Body>
      <S.Footer>
        {Boolean(spansResult[data.id]?.passed) && (
          <S.DotContainer>
            <S.DotContent1>
              <S.DotContent2 $type="success">
                <S.Dot $type="success" />
              </S.DotContent2>
            </S.DotContent1>
          </S.DotContainer>
        )}

        {Boolean(spansResult[data.id]?.failed) && (
          <S.DotContainer>
            <S.DotContent1>
              <S.DotContent2 $type="error">
                <S.Dot $type="error" />
              </S.DotContent2>
            </S.DotContent1>
          </S.DotContainer>
        )}
      </S.Footer>
      <Handle id={id} position={Position.Bottom} style={{bottom: 0, visibility: 'hidden'}} type="source" />
    </S.Container>
  );
};

export default SpanNode;
