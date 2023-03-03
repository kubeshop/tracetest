import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStores} from 'types/DataStore.types';

const SignalFx = () => {
  const baseName = ['dataStore', SupportedDataStores.SignalFX];

  return (
    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="Realm" name={[...baseName, 'realm']} rules={[{required: true, message: 'Realm is required'}]}>
          <Input placeholder="us1" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="Token" name={[...baseName, 'token']} rules={[{required: true, message: 'Token is required'}]}>
          <Input placeholder="Your token" type="password" />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default SignalFx;
