import {ClockCircleOutlined, SettingOutlined, ToolOutlined} from '@ant-design/icons';
import {Handle, Position} from 'react-flow-renderer';

import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import Span from 'models/Span.model';
import * as S from './BaseSpanNode.styled';

interface IProps {
  className: string;
  footer?: React.ReactNode;
  id: string;
  isMatched: boolean;
  isSelected: boolean;
  span: Span;
}

const BaseSpanNode = ({className, footer, id, isMatched, isSelected, span}: IProps) => {
  return (
    <S.Container className={className} data-cy={`trace-node-${span.type}`} $matched={isMatched} $selected={isSelected}>
      <Handle id={id} position={Position.Top} style={{ top: 0, visibility: 'hidden' }} type="target" />

      <S.TopLine $type={span.type} />

      <S.Header>
        <S.BadgeContainer>
          <S.BadgeType count={SemanticGroupNamesToText[span.type]} $hasMargin $type={span.type} />
        </S.BadgeContainer>
        <S.HeaderText>{span.name}</S.HeaderText>
      </S.Header>

      <S.Body>
        <S.Item>
          <SettingOutlined />
          <S.ItemText>
            {span.service} {SpanKindToText[span.kind]}
          </S.ItemText>
        </S.Item>
        {Boolean(span.system) && (
          <S.Item>
            <ToolOutlined />
            <S.ItemText>{span.system}</S.ItemText>
          </S.Item>
        )}
        <S.Item>
          <ClockCircleOutlined />
          <S.ItemText>{span.duration}</S.ItemText>
        </S.Item>
      </S.Body>

      <S.Footer>{footer}</S.Footer>

      <Handle id={id} position={Position.Bottom} style={{ bottom: 0, visibility: 'hidden' }} type="source" />
    </S.Container>
  );
};

export default BaseSpanNode;
