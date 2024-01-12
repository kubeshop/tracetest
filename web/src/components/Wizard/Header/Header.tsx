import {Progress} from 'antd';
import * as S from './Header.styled';

function calculatePercent(activeStep: number, totalSteps: number) {
  return Math.round((activeStep / totalSteps) * 100);
}

interface IProps {
  activeStep: number;
  totalCompleteSteps: number;
}

const Header = ({activeStep, totalCompleteSteps}: IProps) => (
  <S.Container>
    <S.Title>🚀 Setup guide</S.Title>
    <Progress percent={calculatePercent(activeStep, totalCompleteSteps)} showInfo={false} />
    <S.Text>
      {activeStep} of {totalCompleteSteps} steps completed
    </S.Text>
  </S.Container>
);

export default Header;
