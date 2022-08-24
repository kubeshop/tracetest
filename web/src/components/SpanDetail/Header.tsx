import {SettingOutlined, ToolOutlined} from '@ant-design/icons';

import * as SSpanNode from 'components/Diagram/components/DAG/SpanNode.styled';
import {Steps} from 'components/GuidedTour/traceStepList';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import SpanService from 'services/Span.service';
import {TSpan} from 'types/Span.types';
import * as S from './SpanDetail.styled';

interface IProps {
  span?: TSpan;
  totalFailedChecks?: number;
  totalPassedChecks?: number;
}

const Header = ({span, totalFailedChecks, totalPassedChecks}: IProps) => {
  const {kind, name, service, system, type} = SpanService.getSpanInfo(span);

  return (
    <S.Header data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Details)}>
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
        {Boolean(totalPassedChecks) && (
          <S.HeaderCheck>
            <S.HeaderDot $passed />
            {totalPassedChecks}
          </S.HeaderCheck>
        )}
        {Boolean(totalFailedChecks) && (
          <S.HeaderCheck>
            <S.HeaderDot $passed={false} />
            {totalFailedChecks}
          </S.HeaderCheck>
        )}
      </S.Row>
    </S.Header>
  );
};

export default Header;
