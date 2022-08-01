import {SettingOutlined, ToolOutlined} from '@ant-design/icons';

import * as SSpanNode from 'components/Diagram/components/DAG/SpanNode.styled';
import {Steps} from 'components/GuidedTour/traceStepList';
import {SemanticGroupNamesToText} from 'constants/SemanticGroupNames.constants';
import {SpanKindToText} from 'constants/Span.constants';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TSpan} from 'types/Span.types';
import * as S from './SpanDetail.styled';
import SpanService from '../../services/Span.service';

interface IProps {
  span?: TSpan;
}

const SpanHeader = ({span}: IProps) => {
  const {kind, name, service, system, type} = SpanService.getSpanInfo(span);

  return (
    <S.Header data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Details)}>
      <S.Row>
        <SSpanNode.BadgeType count={SemanticGroupNamesToText[type]} $type={type} />
        <S.HeaderTitle level={2}>{name}</S.HeaderTitle>
      </S.Row>
      <S.Row>
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
      </S.Row>
    </S.Header>
  );
};

export default SpanHeader;
