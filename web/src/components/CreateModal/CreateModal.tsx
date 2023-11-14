import {ICreateTestStep} from 'types/Plugins.types';
import CreateTestSteps from './CreateSteps';
import CreateStepFooter from './CreateSteps/CreateStepFooter';
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
  mode: string;
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
  mode,
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
          mode={mode}
        />
      }
    >
      <CreateTestSteps
        isLoading={false}
        onGoTo={onGoTo}
        stepList={stepList}
        selectedKey={activeStep}
        componentFactory={componentFactory}
        mode={mode}
      />
    </S.Modal>
  );
};

export default CreateModal;
