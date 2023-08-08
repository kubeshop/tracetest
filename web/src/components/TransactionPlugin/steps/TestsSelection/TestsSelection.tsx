import {Form} from 'antd';
import {useCallback} from 'react';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {useCreateTransaction} from 'providers/CreateTransaction/CreateTransaction.provider';
import {TDraftTransaction} from 'types/Transaction.types';
import {ComponentNames} from 'constants/Plugins.constants';
import TestsSelectionForm from './TestsSelectionForm';

const TestsSelection = () => {
  const [form] = Form.useForm<TDraftTransaction>();
  const {onNext} = useCreateTransaction();

  const handleSubmit = useCallback(
    ({steps = []}: TDraftTransaction) => {
      onNext({steps});
    },
    [onNext]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title $withSubtitle>Tests</Step.Title>
        <Step.Subtitle>Define the tests and their order of execution for this transaction</Step.Subtitle>
        <Form<TDraftTransaction>
          id={ComponentNames.TestsSelection}
          autoComplete="off"
          data-cy="create-transaction-modal"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={() => null}
        >
          <TestsSelectionForm />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default TestsSelection;
