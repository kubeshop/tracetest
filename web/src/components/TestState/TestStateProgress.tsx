import {Progress} from 'antd';
import * as S from './TestState.styled';

interface IProps {
  label: string;
  percent: number;
}

const TestStateProgress = ({label, percent}: IProps) => (
  <S.Container hasMinWidth>
    <Progress percent={percent} showInfo={false} status="active" strokeLinecap="square" strokeWidth={6} />
    <S.Text>{label}</S.Text>
  </S.Container>
);

export default TestStateProgress;
