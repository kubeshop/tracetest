import {Form} from 'antd';
import {useCallback, useEffect} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {IRpcValues} from 'types/Test.types';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {ComponentNames} from 'constants/Plugins.constants';
import RequestDetailsForm from './RequestDetailsForm';

const RequestDetails = () => {
  const [form] = Form.useForm<IRpcValues>();
  const {onNext, pluginName, draftTest, onIsFormValid} = useCreateTest();
  const onValidate = useValidateTestDraft({pluginName, setIsValid: onIsFormValid});

  const onRefreshData = useCallback(async () => {
    const {url = '', message = '', method = '', auth, metadata = [{}], protoFile} = draftTest as IRpcValues;
    form.setFieldsValue({url, auth, metadata, message, method, protoFile});

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
          id={ComponentNames.RequestDetails}
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
    </Step.Step>
  );
};

export default RequestDetails;
