import {Form} from 'antd';
import {TDraftTest, TDraftTestForm} from 'types/Test.types';
import * as S from 'components/EditTestForm/EditTestForm.styled';
import EditRequestDetails from 'components/EditTestForm/EditRequestDetails/EditRequestDetails';
import {TriggerTypes} from 'constants/Test.constants';
import BasicDetailsForm from '../CreateTestPlugins/Default/steps/BasicDetails/BasicDetailsForm';

export const FORM_ID = 'edit-test';

interface IProps {
  form: TDraftTestForm;
  onSubmit(values: TDraftTest): Promise<void>;
  onValidation(allValues: any, values: TDraftTest): void;
  triggerType: TriggerTypes;
}

const CreateTestForm = ({form, onSubmit, onValidation, triggerType}: IProps) => {
  return (
    <Form<TDraftTest>
      autoComplete="off"
      data-cy="edit-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={onValidation}
    >
      <S.FormContainer>
        <BasicDetailsForm />

        <EditRequestDetails form={form} type={triggerType} />
      </S.FormContainer>
    </Form>
  );
};

export default CreateTestForm;
