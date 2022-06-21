import {Typography} from 'antd';
import {useSpan} from '../../providers/Span/Span.provider';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/traceStepList';
import Highlighted from '../Highlighted';
import * as S from './SpanDetail.styled';

interface IProps {
  title: string;
}

const SpanHeader: React.FC<IProps> = ({title}) => {
  const {searchText} = useSpan();

  return (
    <S.SpanHeader data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Details)}>
      <S.SpanHeaderTitle>Span Details</S.SpanHeaderTitle>
      <Typography.Text type="secondary">
        <Highlighted highlight={searchText} text={title} />
      </Typography.Text>
    </S.SpanHeader>
  );
};

export default SpanHeader;
