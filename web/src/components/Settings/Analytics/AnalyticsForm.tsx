import {Button, Form, Switch} from 'antd';

import {useSettings} from 'providers/Settings/Settings.provider';
import {IDraftSettings} from 'types/Settings.types';
import * as S from '../common/Settings.styled';

const FORM_ID = 'analytics';

const AnalyticsForm = () => {
  const [form] = Form.useForm<IDraftSettings>();
  const {onSubmit} = useSettings();

  return (
    <Form<IDraftSettings> autoComplete="off" form={form} layout="horizontal" name={FORM_ID} onFinish={onSubmit}>
      <Form.Item label="Enable analytics" name="analytics" valuePropName="checked">
        <Switch />
      </Form.Item>
      <S.FooterContainer>
        <Button htmlType="submit" type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default AnalyticsForm;
