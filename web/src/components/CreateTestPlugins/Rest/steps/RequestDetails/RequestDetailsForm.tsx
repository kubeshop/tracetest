import {Form, Input} from 'antd';
import {IHttpValues, TDraftTestForm} from 'types/Test.types';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import RequestDetailsUrlInput from './RequestDetailsUrlInput';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from './RequestDetailsHeadersInput';

export const FORM_ID = 'create-test';

interface IProps {
  form: TDraftTestForm<IHttpValues>;
  isEditing?: boolean;
}

const RequestDetailsForm = ({form, isEditing = false}: IProps) => {
  return (
    <S.InputContainer $isEditing={isEditing}>
      <RequestDetailsUrlInput />
      <RequestDetailsAuthInput form={form} />
      <RequestDetailsHeadersInput />
      <Form.Item className="input-body" data-cy="body" label="Request body" name="body" style={{marginBottom: 0}}>
        <Input.TextArea placeholder="Enter request body text" />
      </Form.Item>
    </S.InputContainer>
  );
};

export default RequestDetailsForm;
