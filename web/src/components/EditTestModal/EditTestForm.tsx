import {Form, Input} from 'antd';
import {TMethod, TTest} from 'types/Test.types';
import EditRequestDetails from './EditRequestDetails/EditRequestDetails';
import {TEditTest, TEditTestForm} from './EditTestModal';
import * as S from './EditTestModal.styled';
import useValidate from './hooks/useValidate';

export const FORM_ID = 'create-test';

interface IProps {
  form: TEditTestForm;
  test: TTest;
  onSubmit(values: TEditTest): Promise<void>;
  onValidation(isValid: boolean): void;
}

const EditTestForm = ({
  form,
  onSubmit,
  onValidation,
  test: {
    name,
    description,
    serviceUnderTest: {request: {method = 'GET' as TMethod, url = '', body = '', headers = []} = {}} = {},
  },
}: IProps) => {
  const handleOnValuesChange = useValidate(onValidation);
  const initialValues = {
    name,
    description,
    method,
    url,
    body,
    headers,
  };

  return (
    <Form<TEditTest>
      autoComplete="off"
      data-cy="edit-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
      initialValues={initialValues}
    >
      <S.FormSection>
        <S.FormSectionTitle>Basic Details</S.FormSectionTitle>
        <Form.Item
          className="input-name"
          data-cy="name"
          label="Name"
          name="name"
          rules={[{required: true, message: 'Please enter a test name'}]}
        >
          <Input placeholder="Enter test name" />
        </Form.Item>
        <Form.Item
          className="input-description"
          data-cy="create-test-description-input"
          label="Description"
          name="description"
          style={{marginBottom: 0}}
          rules={[{required: true, message: 'Please enter a test description'}]}
        >
          <Input.TextArea placeholder="Enter a brief description" />
        </Form.Item>
      </S.FormSection>
      <EditRequestDetails form={form} type="http" />
    </Form>
  );
};

export default EditTestForm;
