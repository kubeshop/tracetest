import {useCallback, useEffect} from 'react';
import {HTTP_METHOD} from 'constants/Common.constants';
import {IHttpValues, TDraftTestForm, THTTPRequest} from 'types/Test.types';
import BasicDetailsForm from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsForm';
import * as S from '../../EditTestForm.styled';
import {IFormProps} from '../EditRequestDetails';

const EditRequestDetailsHttp = ({form, request}: IFormProps) => {
  const typedForm = form as TDraftTestForm<IHttpValues>;

  const getInitialValues = useCallback(async () => {
    const {url, method, headers, body, auth} = request as THTTPRequest;

    typedForm.setFieldsValue({
      url,
      auth,
      method: method as HTTP_METHOD,
      headers,
      body,
    });
  }, [request, typedForm]);

  useEffect(() => {
    getInitialValues();
  }, [getInitialValues]);

  return (
    <S.FormSection>
      <S.FormSectionTitle>Request Details</S.FormSectionTitle>
      <BasicDetailsForm isEditing form={typedForm} />
    </S.FormSection>
  );
};

export default EditRequestDetailsHttp;
