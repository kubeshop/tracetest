import {Form} from 'antd';
import * as S from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails.styled';
import {IHttpValues, TDraftTestForm} from 'types/Test.types';
import {BodyField} from './BodyField/BodyField';
import RequestDetailsAuthInput from './RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from './RequestDetailsHeadersInput';
import RequestDetailsUrlInput from './RequestDetailsUrlInput';

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
      <BodyField setBody={body => form.setFieldsValue({body})} body={Form.useWatch('body', form)} />
    </S.InputContainer>
  );
};

export default RequestDetailsForm;
