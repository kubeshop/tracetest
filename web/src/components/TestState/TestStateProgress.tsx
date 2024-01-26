import {Progress} from 'antd';
import * as S from './TestState.styled';

interface IProps {
  label: string;
  percent: number;
  showInfo?: boolean;
  info?: string;
}

const TestStateProgress = ({label, percent, showInfo, info}: IProps) => (
  <S.Container hasMinWidth>
    <Progress percent={percent} showInfo={false} status="active" strokeLinecap="square" strokeWidth={6} />
    <S.Text>
      {label} {showInfo && info}
    </S.Text>
  </S.Container>
);

export default TestStateProgress;
