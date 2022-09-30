import {ICreateTestStep} from 'types/Plugins.types';
import CreateTestSteps from '../CreateSteps';
import CreateStepFooter from '../CreateSteps/CreateStepFooter';
import * as S from './CreateModal.styled';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  title: string;
  onGoTo(step: string): void;
  stepList: ICreateTestStep[];
  activeStep: string;
  componentFactory({step}: {step: ICreateTestStep}): React.ReactElement;
  onPrev(): void;
  stepNumber: number;
  isLoading: boolean;
  isValid: boolean;
}

const CreateModal = ({
  isOpen,
  onClose,
  title,
  onGoTo,
  stepList,
  activeStep,
  componentFactory,
  onPrev,
  stepNumber,
  isLoading,
  isValid,
}: IProps) => {
  const isLastStep = stepNumber === stepList.length - 1;
  const step = stepList[stepNumber];

  return (
    <S.Modal
      visible={isOpen}
      onCancel={onClose}
      title={<S.Title>{title}</S.Title>}
      width={860}
      footer={
        <CreateStepFooter
          isValid={isValid}
          isLastStep={isLastStep}
          onPrev={onPrev}
          isLoading={isLoading}
          stepNumber={stepNumber}
          step={step}
        />
      }
    >
      <CreateTestSteps
        isLoading={false}
        onGoTo={onGoTo}
        stepList={stepList}
        selectedKey={activeStep}
        componentFactory={componentFactory}
      />
    </S.Modal>
  );
};

export default CreateModal;
