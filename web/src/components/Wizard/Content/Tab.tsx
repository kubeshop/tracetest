import {CheckOutlined} from '@ant-design/icons';
import {isStepCompleted} from 'models/Wizard.model';
import {IWizardStep} from 'types/Wizard.types';
import * as S from './Content.styled';

interface IProps {
  index: number;
  isActive: boolean;
  step: IWizardStep;
}

const Tab = ({index, isActive, step}: IProps) => {
  const {tabComponent: TabComponent} = step;

  return (
    <S.StepTabContainer $isActive={isActive} $isDisabled={!step.isEnabled}>
      {isStepCompleted(step) ? (
        <S.StepTabCheck>
          <CheckOutlined />
        </S.StepTabCheck>
      ) : (
        <S.StepTabNumber>{index}</S.StepTabNumber>
      )}
      <div>
        <S.StepTabTitle $isActive={isActive}>{step.name}</S.StepTabTitle>
        {isStepCompleted(step) && <TabComponent />}
      </div>
    </S.StepTabContainer>
  );
};

export default Tab;
