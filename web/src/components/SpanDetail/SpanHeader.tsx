import {Typography} from 'antd';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/traceStepList';
import * as S from './SpanDetail.styled';

interface IProps {
  title: string;
}

const SpanHeader: React.FC<IProps> = ({title}) => (
  <S.SpanHeader data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Details)}>
    <S.SpanHeaderTitle>Span Details</S.SpanHeaderTitle>
    <Typography.Text type="secondary">{title}</Typography.Text>
  </S.SpanHeader>
);

export default SpanHeader;
