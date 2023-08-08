import {Form} from 'antd';
import {TDraftTestSuite, TDraftTestSuiteForm} from 'types/TestSuite.types';
import BasicDetailsForm from '../CreateTestPlugins/Default/steps/BasicDetails/BasicDetailsForm';
import TestsSelectionForm from '../TestSuitePlugin/steps/TestsSelection/TestsSelectionForm';
import * as S from './EditTestSuiteForm.styled';

export const FORM_ID = 'edit-test';

interface IProps {
  form: TDraftTestSuiteForm;
  testSuite: TDraftTestSuite;
  onSubmit(values: TDraftTestSuite): Promise<void>;
  onValidation(allValues: any, values: TDraftTestSuite): void;
}

const EditTestSuiteForm = ({form, onSubmit, testSuite, onValidation}: IProps) => {
  return (
    <Form<TDraftTestSuite>
      autoComplete="off"
      data-cy="edit-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={onValidation}
      initialValues={testSuite}
    >
      <S.FormContainer>
        <BasicDetailsForm isEditing />

        <TestsSelectionForm />
      </S.FormContainer>
    </Form>
  );
};

export default EditTestSuiteForm;
