import {Button, Form, Switch} from 'antd';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {useCallback} from 'react';
import {ResourceType, TDraftConfig} from 'types/Settings.types';
import SettingService from 'services/Setting.service';
import * as S from '../common/Settings.styled';

const FORM_ID = 'analytics';

const AnalyticsForm = () => {
  const [form] = Form.useForm<TDraftConfig>();
  const {isLoading, onSubmit} = useSettings();
  const {config} = useSettingsValues();

  const handleOnSubmit = useCallback(
    (values: TDraftConfig) => {
      onSubmit([SettingService.getDraftResource(ResourceType.ConfigType, values)]);
    },
    [onSubmit]
  );

  return (
    <Form<TDraftConfig>
      autoComplete="off"
      form={form}
      initialValues={config}
      layout="horizontal"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item hidden name="id" />
      <Form.Item hidden name="name" />

      <S.SwitchContainer>
        <Form.Item name="analyticsEnabled" valuePropName="checked">
          <Switch />
        </Form.Item>
        <S.SwitchLabel htmlFor={`${FORM_ID}_analyticsEnabled`}>Enable analytics</S.SwitchLabel>
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
