import * as S from 'components/SpanDetail/SpanDetail.styled';
import {Button, Typography} from 'antd';
import {PlusOutlined} from '@ant-design/icons';

interface IProps {
  title: string;
  onClick: () => void;
}

export const GenericHttpSpanHeader = ({title, onClick}: IProps): JSX.Element => (
  <S.DetailsHeader>
    <Typography.Text strong>{title}</Typography.Text>
    <Button type="link" icon={<PlusOutlined />} onClick={onClick}>
      Add Assertion
    </Button>
  </S.DetailsHeader>
);
