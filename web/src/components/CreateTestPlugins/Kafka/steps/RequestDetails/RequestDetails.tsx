import {Form} from 'antd';
import {useCallback, useEffect} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {IKafkaValues} from 'types/Test.types';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {ComponentNames} from 'constants/Plugins.constants';
import RequestDetailsForm from './RequestDetailsForm';

const RequestDetails = () => {
  const [form] = Form.useForm<IKafkaValues>();
  const {onNext, pluginName, draftTest, onIsFormValid} = useCreateTest();
  const onValidate = useValidateTestDraft({pluginName, setIsValid: onIsFormValid});

  const onRefreshData = useCallback(async () => {
    const {brokerUrls = [''], topic = '', authentication, sslVerification = false, headers = [{}], messageKey = '', messageValue = ''} = draftTest as IKafkaValues;
    form.setFieldsValue({brokerUrls, topic, authentication, sslVerification, headers, messageKey, messageValue});

    try {
      await form.validateFields();
      onIsFormValid(true);
    } catch (err) {
      onIsFormValid(false);
    }
  }, [draftTest, form, onIsFormValid]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  const handleSubmit = useCallback(
    async (values: IKafkaValues) => {
      onNext(values);
    },
    [onNext]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <Form<IKafkaValues>
          id={ComponentNames.RequestDetails}
          autoComplete="off"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
          initialValues={{
            headers: [{key: '', value: ''}],
          }}
        >
          <RequestDetailsForm form={form} />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default RequestDetails;
