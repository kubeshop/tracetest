import {Form} from 'antd';
import {useCallback, useState} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import CreateStepFooter from 'components/CreateTestSteps/CreateTestStepFooter';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import BasicDetailsForm from './BasicDetailsForm';
import {IDemoTestExample} from '../../../../../constants/Test.constants';

export interface IBasicDetailsValues {
  name: string;
  description: string;
  testSuite: string;
}

const BasicDetails = () => {
  const [isFormValid, setIsFormValid] = useState(false);
  const [selectedDemo, setSelectedDemo] = useState<IDemoTestExample>();
  const [form] = Form.useForm<IBasicDetailsValues>();
  const {onNext} = useCreateTest();

  const handleNext = useCallback(() => {
    form.submit();
  }, [form]);

  const handleSubmit = useCallback(
    ({name, description}: IBasicDetailsValues) => {
      const {url, body, method} = selectedDemo || {};
      onNext({name, description, serviceUnderTest: {triggerSettings: {http: {url, body, method}}}});
    },
    [onNext, selectedDemo]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Provide needed basic information</Step.Title>
        <BasicDetailsForm
          form={form}
          onSubmit={handleSubmit}
          onSelectDemo={setSelectedDemo}
          onValidation={setIsFormValid}
          selectedDemo={selectedDemo}
        />
      </Step.FormContainer>
      <CreateStepFooter isValid={isFormValid} onNext={handleNext} />
    </Step.Step>
  );
};

export default BasicDetails;
