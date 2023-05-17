import {Button, Checkbox, Form, Input, Switch, Typography} from 'antd';
import {useEffect} from 'react';

import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import SettingService from 'services/Setting.service';
import {ResourceType, TDraftLinter} from 'types/Settings.types';
import * as S from '../common/Settings.styled';

const FORM_ID = 'linter';

const LinterForm = () => {
  const [form] = Form.useForm<TDraftLinter>();
  const {isLoading, onSubmit} = useSettings();
  const {linter} = useSettingsValues();
  const standardsEnabled = Form.useWatch(['plugins', 0, 'enabled'], form);
  const securityEnabled = Form.useWatch(['plugins', 1, 'enabled'], form);

  useEffect(() => {
    form.resetFields();
  }, [form, linter]);

  const handleOnSubmit = (values: TDraftLinter) => {
    values.minimumScore = parseInt(String(values?.minimumScore ?? 0), 10);
    onSubmit([SettingService.getDraftResource(ResourceType.LinterType, values)]);
  };

  return (
    <Form<TDraftLinter>
      autoComplete="off"
      form={form}
      initialValues={linter}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item hidden name="id" />
      <Form.Item hidden name="name" />

      <S.SwitchContainer>
        <label htmlFor={`${FORM_ID}_enabled`}>Enable global linter</label>
        <Form.Item name="enabled" valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      <Form.Item
        label="Minimum score"
        name="minimumScore"
        rules={[{required: true, message: 'Minimum score is required'}]}
        wrapperCol={{span: 8}}
      >
        <Input placeholder="100" type="number" />
      </Form.Item>

      <Typography.Title level={3}>Plugins</Typography.Title>
      <S.LinterPluginsContainer>
        <Form.Item hidden name={['plugins', 0, 'name']} />
        <S.SwitchContainer>
          <label htmlFor={`${FORM_ID}_plugins_0_enabled`}>Enable Standards Plugin</label>
          <Form.Item name={['plugins', 0, 'enabled']} valuePropName="checked">
            <Switch />
          </Form.Item>
        </S.SwitchContainer>

        {standardsEnabled && (
          <Form.Item name={['plugins', 0, 'required']} valuePropName="checked" wrapperCol={{span: 8}}>
            <Checkbox>Required</Checkbox>
          </Form.Item>
        )}

        <Form.Item hidden name={['plugins', 1, 'name']} />
        <S.SwitchContainer>
          <label htmlFor={`${FORM_ID}_plugins_1_enabled`}>Enable Security Plugin</label>
          <Form.Item name={['plugins', 1, 'enabled']} valuePropName="checked">
            <Switch />
          </Form.Item>
        </S.SwitchContainer>

        {securityEnabled && (
          <Form.Item name={['plugins', 1, 'required']} valuePropName="checked" wrapperCol={{span: 8}}>
            <Checkbox>Required</Checkbox>
          </Form.Item>
        )}
      </S.LinterPluginsContainer>

      <S.FooterContainer>
        <Button htmlType="submit" loading={isLoading} type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default LinterForm;
