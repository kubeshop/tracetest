import {ICreateTestStep} from 'types/Plugins.types';
import TransactionPlugin from '.';

interface IProps {
  step: ICreateTestStep;
}

const CreateTransactionFactory = ({step: {component}}: IProps) => {
  const Step = TransactionPlugin[component];

  return <Step />;
};

export default CreateTransactionFactory;
