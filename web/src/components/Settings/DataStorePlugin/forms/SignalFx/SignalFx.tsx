import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStores} from 'types/Config.types';

const SignalFx = () => {
  const baseName = ['dataStore', SupportedDataStores.SignalFX];

  return (
    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="Realm" name={[...baseName, 'realm']} rules={[{required: true, message: 'Realm is required'}]}>
          <Input placeholder="Realm" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="Token" name={[...baseName, 'token']} rules={[{required: true, message: 'Token is required'}]}>
          <Input placeholder="Token" type="password" />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default SignalFx;
