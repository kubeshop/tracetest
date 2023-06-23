import {SettingOutlined, ToolOutlined} from '@ant-design/icons';
import {Space, Tooltip} from 'antd';
import AssertionResultChecks from 'components/AssertionResultChecks';
import * as SSpanNode from 'components/Visualization/components/DAG/BaseSpanNode/BaseSpanNode.styled';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import Span from 'models/Span.model';
import SpanService from 'services/Span.service';
import {TAnalyzerError, TTestSpecSummary} from 'types/TestRun.types';
import * as S from './SpanDetail.styled';

interface IProps {
  span?: Span;
  analyzerErrors?: TAnalyzerError[];
  testSpecs?: TTestSpecSummary;
}

const Header = ({span, analyzerErrors, testSpecs}: IProps) => {
  const {kind, name, service, system, type} = SpanService.getSpanInfo(span);

  if (!span) {
    return (
      <S.Header>
        <S.HeaderTitle level={3}>Span Attributes</S.HeaderTitle>
      </S.Header>
    );
  }

  return (
    <S.Header>
      <S.Column>
        <Space>
          <SSpanNode.BadgeType $hasMargin count={SemanticGroupNamesToText[type]} $type={type} />
          {!!analyzerErrors && (
            <Tooltip title="The analyzer found errors in this span">
              <S.LintErrorIcon />
            </Tooltip>
          )}
        </Space>
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
        {!!testSpecs && (
          <AssertionResultChecks failed={testSpecs.failed} passed={testSpecs.passed} styleType="summary" />
        )}
      </S.Row>
    </S.Header>
  );
};

export default Header;
