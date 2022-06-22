import {Form, FormInstance, Input} from 'antd';
import {IRequestDetailsValues} from './RequestDetails';
import * as S from './RequestDetails.styled';
import useValidate from './hooks/useValidate';
import RequestDetailsUrlInput from './RequestDetailsUrlInput';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/CreateTestFormAuthInput';
import RequestDetailsHeadersInput from './RequestDetailsHeadersInput';

export const FORM_ID = 'create-test';

interface IProps {
  form: FormInstance<IRequestDetailsValues>;
  onSubmit(values: IRequestDetailsValues): void;
  onValidation(isValid: boolean): void;
}

const BasicDetailsForm = ({form, onSubmit, onValidation}: IProps) => {
  const handleOnValuesChange = useValidate(onValidation);

  return (
    <Form
      autoComplete="off"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
    >
      <S.GlobalStyle />
      <S.InputContainer>
        <RequestDetailsUrlInput />
        <RequestDetailsAuthInput form={form} />
        <RequestDetailsHeadersInput />
        <Form.Item className="input-body" data-cy="body" label="Request body" name="body" style={{marginBottom: 0}}>
          <Input.TextArea placeholder="Enter request body text" />
        </Form.Item>
      </S.InputContainer>
    </Form>
  );
};

export default BasicDetailsForm;
