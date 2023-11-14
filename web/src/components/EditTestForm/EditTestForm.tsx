import {Form} from 'antd';
import {useMemo} from 'react';
import {TDraftTest, TDraftTestForm} from 'types/Test.types';
import Test from 'models/Test.model';
import TestService from 'services/Test.service';
import FormFactory from 'components/TestPlugins/FormFactory';
import BasicDetailsForm from '../BasicDetailsForm/BasicDetailsForm';
import * as S from './EditTestForm.styled';

export const FORM_ID = 'edit-test';

interface IProps {
  form: TDraftTestForm;
  test: Test;
  onSubmit(values: TDraftTest): Promise<void>;
  onValidation(allValues: any, values: TDraftTest): void;
}

const EditTestForm = ({
  form,
  onSubmit,
  test: {
    trigger: {type},
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
        <BasicDetailsForm />

        <FormFactory type={type} />
      </S.FormContainer>
    </Form>
  );
};

export default EditTestForm;
