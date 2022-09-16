import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import {HTTP_METHOD} from 'constants/Common.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {useCallback, useEffect} from 'react';
import {IHttpValues} from 'types/Test.types';
import RequestDetailsForm from './RequestDetailsForm';

const RequestDetails = () => {
  const [form] = Form.useForm<IHttpValues>();
  const {onNext} = useCreateTest();
  const {draftTest, pluginName} = useCreateTest();
  const {isValid, onValidate, setIsValid} = useValidateTestDraft({pluginName});
  const {url = '', body = '', method = HTTP_METHOD.GET} = draftTest as IHttpValues;

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    (values: IHttpValues) => {
      onNext(values);
    },
    [onNext]
  );

  const onRefreshData = useCallback(async () => {
    form.setFieldsValue({url, body, method: method as HTTP_METHOD});

    try {
      form.validateFields();
      setIsValid(true);
    } catch (err) {
      setIsValid(false);
    }
  }, [body, form, method, setIsValid, url]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Provide additional information</Step.Title>
        <Form<IHttpValues>
          autoComplete="off"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
        >
          <RequestDetailsForm form={form} />
        </Form>
      </Step.FormContainer>
      <CreateStepFooter isValid={isValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default RequestDetails;
