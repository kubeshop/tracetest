import {Form} from 'antd';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import {useCallback, useEffect} from 'react';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {IPostmanValues} from 'types/Test.types';
import UploadCollectionForm from './UploadCollectionForm';

const UploadCollection = () => {
  const [form] = Form.useForm<IPostmanValues>();
  const {onNext, pluginName} = useCreateTest();

  const {isValid, onValidate, setIsValid} = useValidateTestDraft({pluginName});

  const handleOnSubmit = useCallback(
    (values: IPostmanValues) => {
      onNext(values);
    },
    [onNext]
  );

  const url = Form.useWatch('url', form);

  const onValidateUrlChange = useCallback(async () => {
    try {
      await form.validateFields();
      setIsValid(true);
    } catch (err) {
      setIsValid(false);
    }
  }, [form, setIsValid]);

  useEffect(() => {
    onValidateUrlChange();
  }, [url, onValidateUrlChange]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <Form<IPostmanValues>
          autoComplete="off"
          form={form}
          layout="vertical"
          initialValues={{url: '', requests: [], variables: []}}
          onFinish={handleOnSubmit}
          onValuesChange={onValidate}
        >
          <UploadCollectionForm form={form} />
        </Form>
      </Step.FormContainer>
      <CreateStepFooter isValid={isValid} onNext={useCallback(() => form.submit(), [form])} />
    </Step.Step>
  );
};

export default UploadCollection;
