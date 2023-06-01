import {Button, Form, Input, Switch} from 'antd';
import {useEffect} from 'react';

import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import SettingService from 'services/Setting.service';
import {ResourceType, TDraftLinter} from 'types/Settings.types';
// import Plugin from './Plugin';
import * as S from '../common/Settings.styled';

const FORM_ID = 'linter';

const LinterForm = () => {
  const [form] = Form.useForm<TDraftLinter>();
  const {isLoading, onSubmit} = useSettings();
  const {linter} = useSettingsValues();

  useEffect(() => {
    form.resetFields();
  }, [form, linter]);

  const handleOnSubmit = (values: TDraftLinter) => {
    onSubmit([
      SettingService.getDraftResource(ResourceType.AnalyzerType, {
        ...values,
        minimumScore: parseInt(String(values?.minimumScore ?? 0), 10),
        plugins: linter.plugins,
      }),
    ]);
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
        <Form.Item name="enabled" valuePropName="checked">
          <Switch />
        </Form.Item>
        <label htmlFor={`${FORM_ID}_enabled`}>Enable Linter for All Tests</label>
      </S.SwitchContainer>

      <Form.Item
        label="Minimum score"
        name="minimumScore"
        rules={[{required: true, message: 'Minimum score is required'}]}
        wrapperCol={{span: 8}}
      >
        <Input placeholder="0 to 100" type="number" />
      </Form.Item>

      {/* <Typography.Title level={3}>Plugins</Typography.Title>
      <S.LinterPluginsContainer>
        {linter.plugins.map((plugin, index) => (
          <Plugin formId={FORM_ID} index={index} key={plugin.name} plugin={plugin} />
        ))}
      </S.LinterPluginsContainer> */}

      <S.FooterContainer>
        <Button htmlType="submit" loading={isLoading} type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default LinterForm;
