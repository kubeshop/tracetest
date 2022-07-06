import {Form, FormInstance, Input, Select} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import {IRequestDetailsValues} from './RequestDetails';
import * as S from './RequestDetails.styled';
import useValidate from './hooks/useValidate';
import RequestDetailsMetadataInput from './RequestDetailsMetadataInput';
import RequestDetailsFileInput from './RequestDetailsFileInput';

export const FORM_ID = 'create-test';

interface IProps {
  form: FormInstance<IRequestDetailsValues>;
  onSubmit(values: IRequestDetailsValues): void;
  onValidation(isValid: boolean): void;
  methodList: string[];
}

const RequestDetailsForm = ({form, onSubmit, onValidation, methodList}: IProps) => {
  const handleOnValuesChange = useValidate(onValidation);

  return (
    <Form
      autoComplete="off"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
      initialValues={{
        metadata: [{key: '', value: ''}],
      }}
    >
      <S.InputContainer>
        <Form.Item data-cy="protoFile" name="protoFile" label="Upload Protobuf File">
          <RequestDetailsFileInput />
        </Form.Item>
        <Form.Item data-cy="method" label="Select Method" name="method">
          <Select data-cy="method-select">
            {methodList.map(method => (
              <Select.Option data-cy={`rpc-method-${method}`} key={method} value={method}>
                {method}
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
      </S.InputContainer>
      <Step.Title>Provide Additional Information</Step.Title>
      <S.DoubleInputContainer>
        <RequestDetailsUrlInput showMethodSelector={false} />
        <RequestDetailsAuthInput form={form} />
        <RequestDetailsMetadataInput />
        <Form.Item data-cy="message" label="Message" name="message" style={{marginBottom: 0}}>
          <Input.TextArea placeholder="Enter message" />
        </Form.Item>
      </S.DoubleInputContainer>
    </Form>
  );
};

export default RequestDetailsForm;
