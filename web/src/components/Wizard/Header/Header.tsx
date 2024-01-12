import {Progress} from 'antd';
import * as S from './Header.styled';

function calculatePercent(activeStep: number, totalSteps: number) {
  return Math.round((activeStep / totalSteps) * 100);
}

interface IProps {
  activeStep: number;
  totalSteps: number;
}

const Header = ({activeStep, totalSteps}: IProps) => (
  <S.Container>
    <S.Title>🚀 Setup guide</S.Title>
    <Progress percent={calculatePercent(activeStep, totalSteps)} showInfo={false} />
    <S.Text>
      {activeStep} of {totalSteps} steps completed
    </S.Text>
  </S.Container>
);

export default Header;
