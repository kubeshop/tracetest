import {TStepStatus} from '../../../types/Plugins.types';
import * as S from './CreateSteps.styled';

interface IProps {
  text: string;
  status?: TStepStatus;
  isActive: boolean;
}

const CreateTestStepsTab = ({text, status, isActive}: IProps) => {
  return (
    <S.CreateStepsTab>
      <S.StepDot $isActive={isActive} />
      <S.CreateStepsTabTitle $isActive={isActive}>{text}</S.CreateStepsTabTitle>
      {status === 'complete' && <S.StatusIcon />}
    </S.CreateStepsTab>
  );
};

export default CreateTestStepsTab;
