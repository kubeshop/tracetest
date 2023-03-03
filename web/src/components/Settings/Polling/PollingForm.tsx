import {Button, Col, Form, Input, Row} from 'antd';

import {rawToResource, TRawPolling} from 'models/Polling.model';
import {useSettings} from 'providers/Settings/Settings.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from '../common/Settings.styled';

const FORM_ID = 'polling';

const PollingForm = () => {
  const [form] = Form.useForm<TRawPolling>();
  const {isLoading, onSubmit} = useSettings();
  const {polling} = useSettingsValues();

  const handleOnSubmit = (values: TRawPolling) => {
    onSubmit(rawToResource(values));
  };

  return (
    <Form<TRawPolling>
      autoComplete="off"
      form={form}
      initialValues={polling}
      layout="vertical"
      name={FORM_ID}
      onFinish={handleOnSubmit}
    >
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
        <Button htmlType="submit" loading={isLoading} type="primary">
          Save
        </Button>
      </S.FooterContainer>
    </Form>
  );
};

export default PollingForm;
