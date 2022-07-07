import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {useCallback, useState} from 'react';
import {HTTP_METHOD} from '../../../../../constants/Common.constants';
import {TRequest, TRequestAuth} from '../../../../../types/Test.types';
import UploadCollectionForm from './UploadCollectionForm';

export interface IRequestDetailsValues {
  file?: File;
  collectionTest?: string;
  body: string;
  auth: TRequestAuth;
  headers: TRequest['headers'];
  method: HTTP_METHOD;
  name: string;
  url: string;
}

const UploadCollection = () => {
  const [isFormValid, setIsFormValid] = useState(false);
  const [form] = Form.useForm<IRequestDetailsValues>();
  const {onNext} = useCreateTest();
  // const {
  //   draftTest: {serviceUnderTest: {request} = {}},
  // } = useCreateTest();

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    ({file, auth, collectionTest, ...values}: IRequestDetailsValues) => {
      onNext({serviceUnderTest: {request: {...values, auth: {}}}});
    },
    [onNext]
  );

  // const onRefreshData = useCallback(async () => {
  //   const {url = '', body = '', method = HTTP_METHOD.GET} = request || {};
  //
  //   // form.setFieldsValue({url, body, method: method as HTTP_METHOD});
  //
  //   try {
  //     await form.validateFields();
  //     setIsFormValid(true);
  //   } catch (err) {
  //     setIsFormValid(false);
  //   }
  // }, [form, request]);
  //
  // useEffect(() => {
  //   onRefreshData();
  // }, [onRefreshData]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <UploadCollectionForm form={form} onSubmit={handleSubmit} onValidation={setIsFormValid} />
      </Step.FormContainer>
      <CreateStepFooter isValid={isFormValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default UploadCollection;
