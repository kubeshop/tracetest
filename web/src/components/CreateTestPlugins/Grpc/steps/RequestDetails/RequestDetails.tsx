import {Form} from 'antd';
import {useCallback, useEffect} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {IRpcValues} from 'types/Test.types';
import RequestDetailsForm from './RequestDetailsForm';
import useValidateTestDraft from '../../../../../hooks/useValidateTestDraft';

const RequestDetails = () => {
  const [form] = Form.useForm<IRpcValues>();
  const {onNext, pluginName, draftTest} = useCreateTest();
  const {isValid, onValidate, setIsValid} = useValidateTestDraft({pluginName});
  const {url = '', message = '', method = '', auth, metadata = [{}], protoFile} = draftTest as IRpcValues;

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const onRefreshData = useCallback(async () => {
    form.setFieldsValue({url, auth, metadata, message, method, protoFile});

    try {
      await form.validateFields();
      setIsValid(true);
    } catch (err) {
      setIsValid(false);
    }
  }, [auth, message, metadata, method, protoFile, url]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  const handleSubmit = useCallback(
    async (values: IRpcValues) => {
      onNext(values);
    },
    [onNext]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <Form<IRpcValues>
          autoComplete="off"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
          initialValues={{
            metadata: [{key: '', value: ''}],
          }}
        >
          <RequestDetailsForm form={form} />
        </Form>
      </Step.FormContainer>
      <CreateStepFooter isValid={isValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default RequestDetails;
