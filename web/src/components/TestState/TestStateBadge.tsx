import {Badge} from 'antd';
import * as S from './TestState.styled';

interface IProps {
  label: string;
  status: 'success' | 'processing' | 'error' | 'default' | 'warning';
}

const TestStateBadge = ({label, status}: IProps) => (
  <S.Container>
    <Badge status={status} text={label} />
  </S.Container>
);

export default TestStateBadge;
