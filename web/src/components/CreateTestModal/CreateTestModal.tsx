import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {useEffect} from 'react';
import CreateModal from '../CreateModal/CreateModal';
import CreateTestFactory from '../CreateTestPlugins/CreateTestFactory';

interface IProps {
  isOpen: boolean;
  onClose(): void;
}

const CreateTestModal = ({isOpen, onClose}: IProps) => {
  const {stepList, stepNumber, activeStep, onPrev, onReset, isFormValid, onIsFormValid, isLoading} = useCreateTest();

  useEffect(() => {
    if (!isOpen) onReset();
  }, [isOpen, onReset]);

  useEffect(() => {
    const step = stepList[stepNumber];
    onIsFormValid(Boolean(step.isDefaultValid) || step.status === 'complete');
  }, [onIsFormValid, stepList, stepNumber]);

  return isOpen ? (
    <CreateModal
      isOpen
      isValid={isFormValid}
      onClose={onClose}
      title="Create New Test"
      stepList={stepList}
      activeStep={activeStep}
      onGoTo={() => null}
      onPrev={onPrev}
      isLoading={isLoading}
      stepNumber={stepNumber}
      componentFactory={CreateTestFactory}
      mode="CreateTestFactory"
    />
  ) : null;
};

export default CreateTestModal;
