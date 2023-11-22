import {Form, Switch} from 'antd';
import {useCallback} from 'react';

import AllowButton, {Operation} from 'components/AllowButton';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import SettingService from 'services/Setting.service';
import {SupportedDemosFormField, TDraftDemo} from 'types/Settings.types';
import OtelFields from './OtelFields';
import PokeshopFields from './PokeshopFields';
import * as S from '../common/Settings.styled';

const FORM_ID = 'demo';

const DemoForm = () => {
  const [form] = Form.useForm<TDraftDemo>();
  const {isLoading, onSubmit} = useSettings();
  const {demos} = useSettingsValues();
  const pokeshopEnabled = Form.useWatch([SupportedDemosFormField.Pokeshop, 'enabled'], form);
  const otelEnabled = Form.useWatch([SupportedDemosFormField.OpentelemetryStore, 'enabled'], form);

  const handleOnSubmit = useCallback(
    (values: TDraftDemo) => {
      onSubmit(SettingService.getDemoFormValues(values));
    },
    [onSubmit]
  );

  return (
    <Form<TDraftDemo>
      autoComplete="off"
      form={form}
      initialValues={SettingService.getDemoFormInitialValues(demos)}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <Form.Item name={[SupportedDemosFormField.Pokeshop, 'type']} hidden />
      <Form.Item name={[SupportedDemosFormField.Pokeshop, 'id']} hidden />
      <Form.Item name={[SupportedDemosFormField.Pokeshop, 'name']} hidden />

      <S.SwitchContainer>
        <Form.Item name={[SupportedDemosFormField.Pokeshop, 'enabled']} valuePropName="checked" noStyle>
          <Switch />
        </Form.Item>
        <S.SwitchLabel htmlFor={`${FORM_ID}_pokeshop_enabled`}>Enable Pokeshop</S.SwitchLabel>
      </S.SwitchContainer>

      {pokeshopEnabled && <PokeshopFields />}

      <Form.Item name={[SupportedDemosFormField.OpentelemetryStore, 'type']} hidden />
      <Form.Item name={[SupportedDemosFormField.OpentelemetryStore, 'id']} hidden />
      <Form.Item name={[SupportedDemosFormField.OpentelemetryStore, 'name']} hidden />

      <S.SwitchContainer>
        <Form.Item name={[SupportedDemosFormField.OpentelemetryStore, 'enabled']} valuePropName="checked" noStyle>
          <Switch />
        </Form.Item>
        <S.SwitchLabel htmlFor={`${FORM_ID}_opentelemetryStore_enabled`}>
          Enable OpenTelemetry Astronomy Shop Demo
        </S.SwitchLabel>
      </S.SwitchContainer>

      {otelEnabled && <OtelFields />}

      <S.FooterContainer>
        <AllowButton
          operation={Operation.Configure}
          htmlType="submit"
          loading={isLoading}
          type="primary"
          data-cy="demo-form-save-button"
        >
          Save
        </AllowButton>
      </S.FooterContainer>
    </Form>
  );
};

export default DemoForm;
