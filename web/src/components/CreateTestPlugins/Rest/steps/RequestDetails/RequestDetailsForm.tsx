import {Form} from 'antd';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import {IHttpValues, TDraftTestForm} from 'types/Test.types';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import KeyValueListInput from 'components/KeyValueListInput';
import {BodyField} from './BodyField/BodyField';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
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
      <KeyValueListInput name='headers' label='Header list' addButtonLabel='Add Header' keyPlaceholder='Header' valuePlaceholder='Value' initialValue={DEFAULT_HEADERS} />
      <BodyField setBody={body => form.setFieldsValue({body})} body={Form.useWatch('body', form)} />
      <SSLVerification />
    </S.InputContainer>
  );
};

export default RequestDetailsForm;
