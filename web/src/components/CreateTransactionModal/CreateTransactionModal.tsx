import {useEffect} from 'react';
import {useCreateTransaction} from 'providers/CreateTransaction/CreateTransaction.provider';
import CreateModal from '../CreateModal/CreateModal';
import CreateTransactionFactory from '../TransactionPlugin/CreateTransactionFactory';

interface IProps {
  isOpen: boolean;
  onClose(): void;
}

const CreateTransactionModal = ({isOpen, onClose}: IProps) => {
  const {stepList, stepNumber, onPrev, activeStep, onReset, onIsFormValid, isFormValid} = useCreateTransaction();

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
      title="Create Transaction"
      stepList={stepList}
      activeStep={activeStep}
      onGoTo={() => null}
      onPrev={onPrev}
      isLoading={false}
      stepNumber={stepNumber}
      componentFactory={CreateTransactionFactory}
      mode="CreateTransactionFactory"
    />
  ) : null;
};

export default CreateTransactionModal;
