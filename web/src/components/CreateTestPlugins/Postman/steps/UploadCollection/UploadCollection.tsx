import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import {HTTP_METHOD} from 'constants/Common.constants';
import {VariableDefinition} from 'postman-collection';
import {useCallback, useState} from 'react';
import {RequestDefinitionExtended} from 'services/PostmanService.service';
import {THTTPRequest, TRequestAuth} from 'types/Test.types';
import Validator from 'utils/Validator';
import {useOnSubmitCallback} from './hooks/useOnSubmitCallback';
import {useValidateFormEffect} from './hooks/useValidateFormEffect';
import UploadCollectionForm from './UploadCollectionForm';

export interface IUploadCollectionValues {
  collectionFile?: File;
  envFile?: File;
  collectionTest?: string;
  requests: RequestDefinitionExtended[];
  variables: VariableDefinition[];
  body: string;
  auth: TRequestAuth;
  headers: THTTPRequest['headers'];
  method: HTTP_METHOD;
  url: string;
}

const UploadCollection = () => {
  const [transientUrl, setTransientUrl] = useState('');
  const [form] = Form.useForm<IUploadCollectionValues>();
  const [isFormValid, setIsFormValid] = useValidateFormEffect(form);
  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <UploadCollectionForm
          setTransientUrl={setTransientUrl}
          form={form}
          onSubmit={useOnSubmitCallback()}
          onValidation={setIsFormValid}
        />
      </Step.FormContainer>
      <CreateStepFooter
        isValid={isFormValid && Validator.url(transientUrl)}
        onNext={useCallback(() => form.submit(), [form])}
      />
    </Step.Step>
  );
};

export default UploadCollection;
