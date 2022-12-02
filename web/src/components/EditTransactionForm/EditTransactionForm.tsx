import {Form} from 'antd';
import {TDraftTransaction, TDraftTransactionForm} from 'types/Transaction.types';
import BasicDetailsForm from '../CreateTestPlugins/Default/steps/BasicDetails/BasicDetailsForm';
import TestsSelectionForm from '../TransactionPlugin/steps/TestsSelection/TestsSelectionForm';
import * as S from './EditTransactionForm.styled';

export const FORM_ID = 'edit-test';

interface IProps {
  form: TDraftTransactionForm;
  transaction: TDraftTransaction;
  onSubmit(values: TDraftTransaction): Promise<void>;
  onValidation(allValues: any, values: TDraftTransaction): void;
}

const EditTransactionForm = ({form, onSubmit, transaction, onValidation}: IProps) => {
  return (
    <Form<TDraftTransaction>
      autoComplete="off"
      data-cy="edit-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={onValidation}
      initialValues={transaction}
    >
      <S.FormContainer>
        <BasicDetailsForm isEditing />

        <TestsSelectionForm />
      </S.FormContainer>
    </Form>
  );
};

export default EditTransactionForm;
