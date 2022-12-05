import {SettingOutlined, ToolOutlined} from '@ant-design/icons';
import {useMemo} from 'react';
import * as SSpanNode from 'components/Visualization/components/DAG/SpanNode.styled';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import SpanService from 'services/Span.service';
import {TSpan} from 'types/Span.types';
import {TResultAssertions} from 'types/Assertion.types';
import * as S from './SpanDetail.styled';
import AssertionResultChecks from '../AssertionResultChecks/AssertionResultChecks';

interface IProps {
  span?: TSpan;
  assertions?: TResultAssertions;
}

const Header = ({span, assertions = {}}: IProps) => {
  const {kind, name, service, system, type} = SpanService.getSpanInfo(span);
  const {failed, passed} = useMemo(() => SpanService.getAssertionResultSummary(assertions), [assertions]);

  return (
    <S.Header>
      <S.Column>
        <SSpanNode.BadgeType $hasMargin count={SemanticGroupNamesToText[type]} $type={type} />
        <S.HeaderTitle level={3}>{name}</S.HeaderTitle>
      </S.Column>
      <S.Column>
        <S.HeaderItem>
          <SettingOutlined />
          <S.HeaderItemText>{`${service} ${SpanKindToText[kind]}`}</S.HeaderItemText>
        </S.HeaderItem>
        {Boolean(system) && (
          <S.HeaderItem>
            <ToolOutlined />
            <S.HeaderItemText>{system}</S.HeaderItemText>
          </S.HeaderItem>
        )}
      </S.Column>
      <S.Row>
        <AssertionResultChecks failed={failed} passed={passed} styleType="summary" />
      </S.Row>
    </S.Header>
  );
};

export default Header;
