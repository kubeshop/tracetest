import {useEffect} from 'react';
import {useCreateTestSuite} from 'providers/CreateTestSuite';
import CreateModal from '../CreateModal/CreateModal';
import CreateTestSuiteFactory from '../TestSuitePlugin/CreateTestSuiteFactory';

interface IProps {
  isOpen: boolean;
  onClose(): void;
}

const CreateTestSuiteModal = ({isOpen, onClose}: IProps) => {
  const {stepList, stepNumber, onPrev, activeStep, onReset, onIsFormValid, isFormValid} = useCreateTestSuite();

  useEffect(() => {
    if (!isOpen) onReset();
  }, [isOpen, onReset]);

  useEffect(() => {
    const step = stepList[stepNumber];
    onIsFormValid(Boolean(step.isDefaultValid) || step.status === 'complete');
  }, [onIsFormValid, stepList, stepNumber]);

  return isOpen ? (
    <CreateModal
      isValid={isFormValid}
      isOpen
      onClose={onClose}
      title="Create Test Suite"
      stepList={stepList}
      activeStep={activeStep}
      onGoTo={() => null}
      onPrev={onPrev}
      isLoading={false}
      stepNumber={stepNumber}
      componentFactory={CreateTestSuiteFactory}
      mode="CreateTestSuiteFactory"
    />
  ) : null;
};

export default CreateTestSuiteModal;
