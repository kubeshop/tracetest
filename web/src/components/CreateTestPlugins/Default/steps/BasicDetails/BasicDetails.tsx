import {Form} from 'antd';
import {useCallback} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {ComponentNames} from 'constants/Plugins.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {IBasicValues} from 'types/Test.types';
import BasicDetailsForm from './BasicDetailsForm';

const BasicDetails = () => {
  const [form] = Form.useForm<IBasicValues>();

  const {
    plugin: {name: pluginName},
    onNext,
    onIsFormValid,
  } = useCreateTest();
  const onValidate = useValidateTestDraft({pluginName, isBasicDetails: true, setIsValid: onIsFormValid});

  const handleSubmit = useCallback(
    (values: IBasicValues) => {
      const {name, description} = values;
      onNext({name, description});
    },
    [onNext]
  );

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Provide needed basic information</Step.Title>
        <Form<IBasicValues>
          id={ComponentNames.BasicDetails}
          autoComplete="off"
          data-cy="create-test-modal"
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          onValuesChange={onValidate}
        >
          <BasicDetailsForm />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default BasicDetails;
