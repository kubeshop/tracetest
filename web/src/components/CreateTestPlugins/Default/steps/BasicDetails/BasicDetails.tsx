import {Form} from 'antd';
import {useCallback, useState} from 'react';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {ComponentNames} from 'constants/Plugins.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {isDemoEnabled} from 'models/Demo.model';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {IBasicValues, TDraftTest} from 'types/Test.types';
import BasicDetailsForm from './BasicDetailsForm';

const BasicDetails = () => {
  const [selectedDemo, setSelectedDemo] = useState<TDraftTest>();
  const [form] = Form.useForm<IBasicValues>();

  const {
    plugin: {name: pluginName, demoList},
    onNext,
    onIsFormValid,
  } = useCreateTest();
  const onValidate = useValidateTestDraft({pluginName, isBasicDetails: true, setIsValid: onIsFormValid});
  const {demo: demoSettings} = useSettingsValues();

  const handleSelectDemo = useCallback(
    (demo: TDraftTest) => {
      form.setFieldsValue(demo);

      onIsFormValid(true);
      setSelectedDemo(demo);
    },
    [form, onIsFormValid]
  );

  const handleSubmit = useCallback(
    (values: IBasicValues) => {
      const {name, description} = values;
      onNext({...(selectedDemo || {}), name, description});
    },
    [onNext, selectedDemo]
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
          <BasicDetailsForm
            onSelectDemo={handleSelectDemo}
            selectedDemo={selectedDemo}
            demoList={demoList}
            isDemoEnabled={isDemoEnabled(demoSettings)}
          />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default BasicDetails;
