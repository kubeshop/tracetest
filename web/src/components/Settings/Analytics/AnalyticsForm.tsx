import {Button, Form, Switch} from 'antd';

import {rawToResource, TRawConfig} from 'models/Config.model';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from '../common/Settings.styled';

const FORM_ID = 'analytics';

const AnalyticsForm = () => {
  const [form] = Form.useForm<TRawConfig>();
  const {isLoading, onSubmit} = useSettings();
  const {config} = useSettingsValues();

  const handleOnSubmit = (values: TRawConfig) => {
    onSubmit(rawToResource(values));
  };

  return (
    <Form<TRawConfig>
      autoComplete="off"
      form={form}
      initialValues={config}
      layout="horizontal"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item hidden name="id" />

      <S.SwitchContainer>
        <label htmlFor={`${FORM_ID}_analyticsEnabled`}>Enable analytics</label>
        <Form.Item name="analyticsEnabled" valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      <S.FooterContainer>
        <Button htmlType="submit" loading={isLoading} type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default AnalyticsForm;
