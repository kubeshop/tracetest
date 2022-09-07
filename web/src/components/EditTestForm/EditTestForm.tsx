import {Form} from 'antd';
import {useMemo} from 'react';
import {TDraftTest, TDraftTestForm, TTest} from 'types/Test.types';
import TestService from '../../services/Test.service';
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
    trigger: {request, type},
  },
  test,
  onValidation,
}: IProps) => {
  const initialValues = useMemo(() => TestService.getInitialValues(test), [test]);

  return (
    <Form<TDraftTest>
      autoComplete="off"
      data-cy="edit-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={onValidation}
      initialValues={initialValues}
    >
      <S.FormContainer>
        <BasicDetailsForm isEditing />

        <EditRequestDetails form={form} type={type} request={request} />
      </S.FormContainer>
    </Form>
  );
};

export default EditTestForm;
