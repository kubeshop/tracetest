import {Form} from 'antd';
import {useCallback} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {IRpcValues} from 'types/Test.types';
import RequestDetailsForm from './RequestDetailsForm';
import useValidateTestDraft from '../../../../../hooks/useValidateTestDraft';

const RequestDetails = () => {
  const [form] = Form.useForm<IRpcValues>();
  const {
    onNext,
    pluginName,
  } = useCreateTest();
  const {isValid, onValidate} = useValidateTestDraft({pluginName});

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    async ({protoFile, message, metadata, method, auth, url}: IRpcValues) => {
      onNext({
        url,
        message,
        auth,
        method,
        metadata,
        protoFile,
      });
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
