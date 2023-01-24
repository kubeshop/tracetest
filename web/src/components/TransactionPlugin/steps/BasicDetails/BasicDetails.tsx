import {Form} from 'antd';
import {useCallback} from 'react';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {useCreateTransaction} from 'providers/CreateTransaction/CreateTransaction.provider';
import {TDraftTransaction} from 'types/Transaction.types';
import BasicDetailsForm from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetailsForm';
import {ComponentNames} from 'constants/Plugins.constants';

const BasicDetails = () => {
  const [form] = Form.useForm<TDraftTransaction>();
  const {onNext, onIsFormValid} = useCreateTransaction();

  const handleSubmit = useCallback(
    ({name, description}: TDraftTransaction) => {
      onNext({name, description});
    },
    [onNext]
  );

  const onValidate = useCallback(
    async (changedValues, draft: TDraftTransaction) => {
      onIsFormValid(!!draft.name);
    },
    [onIsFormValid]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title $withSubtitle>Name your Transaction</Step.Title>
        <Step.Subtitle>
          Create a transaction which will run several tests in sequence. You can set outputs into variables and use
          these variables later in the transaction in other tests.{' '}
        </Step.Subtitle>
        <Form<TDraftTransaction>
          id={ComponentNames.BasicDetails}
          autoComplete="off"
          data-cy="create-test-modal"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
        >
          <BasicDetailsForm onSelectDemo={() => null} selectedDemo={undefined} demoList={[]} />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default BasicDetails;
