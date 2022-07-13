import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import {VariableDefinition} from 'postman-collection';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {useCallback, useEffect, useState} from 'react';
import {HTTP_METHOD} from 'constants/Common.constants';
import {THTTPRequest, TRequestAuth} from 'types/Test.types';
import Validator from 'utils/Validator';
import {RequestDefinitionExtended} from './hooks/getRequestsFromCollection';
import UploadCollectionForm from './UploadCollectionForm';

export interface IRequestDetailsValues {
  collectionFile?: File;
  envFile?: File;
  collectionTest?: string;
  requests: RequestDefinitionExtended[];
  variables: VariableDefinition[];
  body: string;
  auth: TRequestAuth;
  headers: THTTPRequest['headers'];
  method: HTTP_METHOD;
  name: string;
  url: string;
}

const UploadCollection = () => {
  const [transientUrl, setTransientUrl] = useState('');
  const [isFormValid, setIsFormValid] = useState(false);
  const [form] = Form.useForm<IRequestDetailsValues>();
  const {onNext} = useCreateTest();

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    ({collectionFile, envFile, collectionTest, requests, variables, ...values}: IRequestDetailsValues) => {
      // eslint-disable-next-line no-console
      console.log(collectionFile, envFile, collectionTest, requests, variables);
      onNext({serviceUnderTest: {triggerSettings: {http: values}}});
    },
    [onNext]
  );

  const onRefreshData = useCallback(async () => {
    try {
      await form.validateFields();
      setIsFormValid(true);
    } catch (err) {
      setIsFormValid(false);
    }
  }, [form]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <UploadCollectionForm
          setTransientUrl={setTransientUrl}
          form={form}
          onSubmit={handleSubmit}
          onValidation={setIsFormValid}
        />
      </Step.FormContainer>
      <CreateStepFooter isValid={isFormValid && Validator.url(transientUrl)} onNext={handleNext} />
    </Step.Step>
  );
};

export default UploadCollection;
