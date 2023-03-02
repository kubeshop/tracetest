import {Button, Col, Form, Input, Row} from 'antd';

import {useSettings} from 'providers/Settings/Settings.provider';
import {IDraftSettings} from 'types/Settings.types';
import * as S from '../common/Settings.styled';

const FORM_ID = 'polling';

const PollingForm = () => {
  const [form] = Form.useForm<IDraftSettings>();
  const {onSubmit} = useSettings();

  return (
    <Form<IDraftSettings> autoComplete="off" form={form} layout="vertical" name={FORM_ID} onFinish={onSubmit}>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Max wait time for trace"
            name="maxWaitTimeForTrace"
            rules={[{required: true, message: 'Max wait time for trace is required'}]}
          >
            <Input placeholder="10s" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            label="Retry delay"
            name="retryDelay"
            rules={[{required: true, message: 'Retry delay is required'}]}
          >
            <Input placeholder="500ms" />
          </Form.Item>
        </Col>
      </Row>

      <S.FooterContainer>
        <Button htmlType="submit" type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default PollingForm;
