import {Button, Form, Switch} from 'antd';

import {IDraftSettings} from 'types/Settings.types';
import * as S from '../common/Settings.styled';
import OtelFields from './OtelFields';
import PokeshopFields from './PokeshopFields';

const FORM_ID = 'demo';

const DemoForm = () => {
  const [form] = Form.useForm<IDraftSettings>();
  const pokeshopEnabled = Form.useWatch(['demo', 'pokeshopEnabled'], form);
  const otelEnabled = Form.useWatch(['demo', 'otelEnabled'], form);

  return (
    <Form<IDraftSettings> autoComplete="off" form={form} layout="vertical" name={FORM_ID} onFinish={() => {}}>
      <S.SwitchContainer>
        <label htmlFor={`${FORM_ID}_demo_pokeshopEnabled`}>Enable Pokeshop</label>
        <Form.Item name={['demo', 'pokeshopEnabled']} valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      {pokeshopEnabled && <PokeshopFields />}

      <S.SwitchContainer>
        <label htmlFor={`${FORM_ID}_demo_otelEnabled`}>Enable OpenTelemetry Astronomy Shop Demo</label>
        <Form.Item name={['demo', 'otelEnabled']} valuePropName="checked">
          <Switch />
        </Form.Item>
      </S.SwitchContainer>

      {otelEnabled && <OtelFields />}

      <S.FooterContainer>
        <Button htmlType="submit" type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default DemoForm;
