import {useCallback, useEffect} from 'react';
import RequestDetailsForm from 'components/CreateTestPlugins/Grpc/steps/RequestDetails/RequestDetailsForm';
import {IRpcValues, TDraftTestForm, TGRPCRequest} from 'types/Test.types';
import * as S from '../../EditTestForm.styled';
import {IFormProps} from '../EditRequestDetails';

const EditRequestDetailGrpc = ({form, request}: IFormProps) => {
  const getInitialValues = useCallback(async () => {
    const {address: url, method, metadata, request: message, auth, protobufFile} = request as TGRPCRequest;
    const protoFile = new File([protobufFile], 'file.proto');
    const typedForm = form as TDraftTestForm<IRpcValues>;

    typedForm.setFieldsValue({
      url,
      auth,
      method,
      message,
      metadata,
      protoFile,
    });
  }, [form, request]);

  useEffect(() => {
    getInitialValues();
  }, [getInitialValues]);

  return (
    <S.FormSection>
      <S.FormSectionTitle>Request Details</S.FormSectionTitle>
      <RequestDetailsForm form={form} />
    </S.FormSection>
  );
};

export default EditRequestDetailGrpc;
