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
    <div>
      <S.Title>🚀 Get Started</S.Title>
      <S.Text>Use this guide to get your environment up and running.</S.Text>
    </div>
    <S.ProgressContainer>
      <Progress percent={calculatePercent(activeStep, totalCompleteSteps)} width={300} showInfo={false} />
      <S.Text>
        {activeStep} of {totalCompleteSteps} steps completed
      </S.Text>
    </S.ProgressContainer>
  </S.Container>
);

export default Header;
