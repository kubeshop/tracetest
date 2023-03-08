import {Button, Form, Switch} from 'antd';

import {rawToResource, TRawDemo} from 'models/Demo.model';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import OtelFields from './OtelFields';
import PokeshopFields from './PokeshopFields';
import * as S from '../common/Settings.styled';

const FORM_ID = 'demo';

const DemoForm = () => {
  const [form] = Form.useForm<TRawDemo>();
  const {isLoading, onSubmit} = useSettings();
  const {demo} = useSettingsValues();
  const pokeshopEnabled = Form.useWatch('pokeshopEnabled', form);
  const otelEnabled = Form.useWatch('otelEnabled', form);

  const handleOnSubmit = (values: TRawDemo) => {
    onSubmit(rawToResource(values));
  };

  return (
    <Form<TRawDemo>
      autoComplete="off"
      form={form}
      initialValues={demo}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
      <S.SwitchContainer>
        <label htmlFor={`${FORM_ID}_pokeshopEnabled`}>Enable Pokeshop</label>
        <Form.Item name="pokeshopEnabled" valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      {pokeshopEnabled && <PokeshopFields />}

      <S.SwitchContainer>
        <label htmlFor={`${FORM_ID}_otelEnabled`}>Enable OpenTelemetry Astronomy Shop Demo</label>
        <Form.Item name="otelEnabled" valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      {otelEnabled && <OtelFields />}

      <S.FooterContainer>
        <Button htmlType="submit" loading={isLoading} type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default DemoForm;
