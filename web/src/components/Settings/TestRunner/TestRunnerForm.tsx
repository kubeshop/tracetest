import {useCallback, useEffect} from 'react';
import {Button, Form} from 'antd';
import {ResourceType, TDraftTestRunner} from 'types/Settings.types';
import {useSettings} from 'providers/Settings/Settings.provider';
import SettingService from 'services/Setting.service';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from '../common/Settings.styled';
import RequiredGatesInput from './RequiredGatesInput';

const FORM_ID = 'testRunner';

const TestRunnerForm = () => {
  const [form] = Form.useForm<TDraftTestRunner>();
  const {isLoading, onSubmit} = useSettings();
  const {testRunner} = useSettingsValues();

  const handleOnSubmit = useCallback(
    (values: TDraftTestRunner) => {
      onSubmit([SettingService.getDraftResource(ResourceType.TestRunnerType, values)]);
    },
    [onSubmit]
  );

  useEffect(() => {
    form.setFieldsValue(testRunner);
  }, [form, testRunner]);

  return (
    <Form<TDraftTestRunner>
      autoComplete="off"
      form={form}
      initialValues={testRunner}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item hidden name="id" />
      <Form.Item hidden name="name" />

      <Form.Item name="requiredGates">
        <RequiredGatesInput />
      </Form.Item>

      <S.FooterContainer>
        <Button htmlType="submit" loading={isLoading} type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default TestRunnerForm;
