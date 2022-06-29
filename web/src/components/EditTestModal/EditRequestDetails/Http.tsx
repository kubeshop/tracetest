import {Form, FormInstance, Input} from 'antd';
import RequestDetailsAuthInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import RequestDetailsUrlInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsUrlInput';
import {TEditTest} from '../EditTestModal';
import * as S from '../EditTestModal.styled';

interface IProps {
  form: FormInstance<TEditTest>;
}

const EditRequestDetailsHttp = ({form}: IProps) => {
  return (
    <S.FormSection>
      <S.FormSectionTitle>Request Details</S.FormSectionTitle>
      <RequestDetailsUrlInput />
      <RequestDetailsHeadersInput />
      <Form.Item className="input-body" data-cy="body" label="Request body" name="body">
        <Input.TextArea placeholder="Enter request body text" />
      </Form.Item>
      <RequestDetailsAuthInput form={form} />
    </S.FormSection>
  );
};

export default EditRequestDetailsHttp;
