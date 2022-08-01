import {Form} from 'antd';
import {useCallback, useState} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {IBasicValues, TDraftTest} from 'types/Test.types';
import BasicDetailsForm from './BasicDetailsForm';

const BasicDetails = () => {
  const [selectedDemo, setSelectedDemo] = useState<TDraftTest>();
  const [form] = Form.useForm<IBasicValues>();
  const {onNext} = useCreateTest();

  const {
    plugin: {name: pluginName, demoList},
  } = useCreateTest();
  const {setIsValid, isValid, onValidate} = useValidateTestDraft({pluginName, isBasicDetails: true});

  const handleSelectDemo = useCallback(
    (demo: TDraftTest) => {
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
      onNext({...(selectedDemo || {}), name, description});
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
            demoList={demoList}
          />
        </Form>
      </Step.FormContainer>
      <CreateStepFooter isValid={isValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default BasicDetails;
