import {Form} from 'antd';
import {useCallback, useState} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {IDemoTestExample} from 'constants/Test.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {SupportedPlugins} from 'constants/Plugins.constants';
import {IBasicValues} from 'types/Test.types';
import BasicDetailsForm from './BasicDetailsForm';

const BasicDetails = () => {
  const [selectedDemo, setSelectedDemo] = useState<IDemoTestExample>();
  const [form] = Form.useForm<IBasicValues>();
  const {onNext} = useCreateTest();

  const {
    plugin: {name: pluginName},
  } = useCreateTest();
  const {setIsValid, isValid, onValidate} = useValidateTestDraft({pluginName, isBasicDetails: true});

  const handleSelectDemo = useCallback(
    (demo: IDemoTestExample) => {
      const {name, description} = demo;

      form.setFieldsValue({
        name,
        description,
      });

      setIsValid(true);
      setSelectedDemo(demo);
    },
    [form, setIsValid]
  );

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    ({name, description}: IBasicValues) => {
      const {url, body, method} = selectedDemo || {};
      onNext({name, description, url, body, method});
    },
    [onNext, selectedDemo]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Provide needed basic information</Step.Title>
        <Form<IBasicValues>
          autoComplete="off"
          data-cy="create-test-modal"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
        >
          <BasicDetailsForm
            onSelectDemo={handleSelectDemo}
            selectedDemo={selectedDemo}
            showDemo={pluginName === SupportedPlugins.REST}
          />
        </Form>
      </Step.FormContainer>
      <CreateStepFooter isValid={isValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default BasicDetails;
