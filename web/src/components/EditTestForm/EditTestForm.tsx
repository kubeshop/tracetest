import {Form} from 'antd';
import {TDraftTest, TDraftTestForm, TTest} from 'types/Test.types';
import BasicDetailsForm from '../CreateTestPlugins/Default/steps/BasicDetails/BasicDetailsForm';
import EditRequestDetails from './EditRequestDetails/EditRequestDetails';
import * as S from './EditTestForm.styled';

export const FORM_ID = 'edit-test';

interface IProps {
  form: TDraftTestForm;
  test: TTest;
  onSubmit(values: TDraftTest): Promise<void>;
  onValidation(allValues: any, values: TDraftTest): void;
}

const EditTestForm = ({
  form,
  onSubmit,
  test: {
    name,
    description,
    trigger: {request, type},
  },
  onValidation,
}: IProps) => {
  return (
    <Form<TDraftTest>
      autoComplete="off"
      data-cy="edit-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={onValidation}
      initialValues={{
        name,
        description,
        type,
      }}
    >
      <S.FormSection>
        <S.FormSectionTitle>Basic Details</S.FormSectionTitle>
        <BasicDetailsForm isEditing />
      </S.FormSection>

      <EditRequestDetails form={form} type={type} request={request} />
    </Form>
  );
};

export default EditTestForm;
