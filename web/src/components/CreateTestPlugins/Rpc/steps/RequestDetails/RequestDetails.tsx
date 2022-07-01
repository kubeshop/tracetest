import {Form} from 'antd';
import {useCallback, useEffect, useState} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import RpcService from 'services/Rpc.service';
import {TRequestAuth, TRequest} from 'types/Test.types';
import RequestDetailsForm from './RequestDetailsForm';

export interface IRequestDetailsValues {
  message: string;
  auth: TRequestAuth;
  metadata: TRequest['headers'];
  url: string;
  method: string;
  protoFile: File;
}

const RequestDetails = () => {
  const [isFormValid, setIsFormValid] = useState(false);
  const [methodList, setMethodList] = useState<string[]>([]);
  const [form] = Form.useForm<IRequestDetailsValues>();
  const {onNext} = useCreateTest();

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    (values: IRequestDetailsValues) => {
      console.log('@@values', values);
      onNext();
    },
    [onNext]
  );

  const protoFile = Form.useWatch('protoFile', form);

  const getMethodList = useCallback(async () => {
    if (protoFile) {
      const list = await RpcService.getMethodList(protoFile);

      setMethodList(list);
    } else {
      setMethodList([]);
      form.setFieldsValue({
        method: '',
      });
    }
  }, [form, protoFile]);

  useEffect(() => {
    getMethodList();
  }, [getMethodList]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <RequestDetailsForm form={form} methodList={methodList} onSubmit={handleSubmit} onValidation={setIsFormValid} />
      </Step.FormContainer>
      <CreateStepFooter isValid={isFormValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default RequestDetails;
