import {Typography} from 'antd';
import * as S from './SpanDetail.styled';

interface IProps {
  title: string;
}

const SpanHeader: React.FC<IProps> = ({title}) => (
  <S.SpanHeader>
    <S.SpanHeaderTitle>Span Details</S.SpanHeaderTitle>
    <Typography.Text type="secondary">{title}</Typography.Text>
  </S.SpanHeader>
);

export default SpanHeader;
