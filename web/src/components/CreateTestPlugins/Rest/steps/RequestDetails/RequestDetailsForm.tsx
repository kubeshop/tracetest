import {Form} from 'antd';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import {IHttpValues, TDraftTestForm} from 'types/Test.types';
import {BodyField} from './BodyField/BodyField';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from './RequestDetailsHeadersInput';
import RequestDetailsUrlInput from './RequestDetailsUrlInput';
import SSLVerification from './SSLVerification';

export const FORM_ID = 'create-test';

interface IProps {
  form: TDraftTestForm<IHttpValues>;
}

const RequestDetailsForm = ({form}: IProps) => {
  return (
    <S.InputContainer>
      <RequestDetailsUrlInput />
      <RequestDetailsAuthInput />
      <RequestDetailsHeadersInput />
      <BodyField setBody={body => form.setFieldsValue({body})} body={Form.useWatch('body', form)} />
      <SSLVerification />
    </S.InputContainer>
  );
};

export default RequestDetailsForm;
