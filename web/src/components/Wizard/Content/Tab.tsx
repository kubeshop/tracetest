import {CheckOutlined} from '@ant-design/icons';
import {isStepCompleted} from 'models/Wizard.model';
import {IWizardStep} from 'types/Wizard.types';
import * as S from './Content.styled';

interface IProps {
  index: number;
  isActive: boolean;
  step: IWizardStep;
}

const Tab = ({index, isActive, step}: IProps) => (
  <S.StepTabContainer $isActive={isActive}>
    {isStepCompleted(step) ? (
      <S.StepTabCheck>
        <CheckOutlined />
      </S.StepTabCheck>
    ) : (
      <S.StepTabNumber>{index}</S.StepTabNumber>
    )}
    <S.StepTabTitle $isActive={isActive}>{step.name}</S.StepTabTitle>
  </S.StepTabContainer>
);

export default Tab;
