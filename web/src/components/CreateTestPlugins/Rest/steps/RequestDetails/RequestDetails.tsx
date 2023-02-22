import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {HTTP_METHOD} from 'constants/Common.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {ComponentNames} from 'constants/Plugins.constants';
import {useCallback, useEffect} from 'react';
import {IHttpValues} from 'types/Test.types';
import {DEFAULT_HEADERS} from 'constants/Test.constants';
import RequestDetailsForm from './RequestDetailsForm';

const RequestDetails = () => {
  const {onNext, draftTest, pluginName, onIsFormValid} = useCreateTest();
  const [form] = Form.useForm<IHttpValues>();
  const onValidate = useValidateTestDraft({pluginName, setIsValid: onIsFormValid});

  const handleSubmit = useCallback(
    (values: IHttpValues) => {
      onNext(values);
    },
    [onNext]
  );

  const onRefreshData = useCallback(async () => {
    const {url = '', body = '', method = HTTP_METHOD.GET, headers = DEFAULT_HEADERS} = draftTest as IHttpValues;
    form.setFieldsValue({url, body, method: method as HTTP_METHOD, headers});

    onValidate(null, form.getFieldsValue());
  }, [draftTest, form, onValidate]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Provide additional information</Step.Title>
        <Form<IHttpValues>
          id={ComponentNames.RequestDetails}
          autoComplete="off"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
        >
          <RequestDetailsForm form={form} />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default RequestDetails;
