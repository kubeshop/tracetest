import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStores} from 'types/Config.types';

const AwsXRay = () => {
  const baseName = ['dataStore', SupportedDataStores.AWSXRay];

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Access Key Id"
            name={[...baseName, 'accessKeyId']}
            rules={[{required: true, message: 'Access Key Id is required'}]}
          >
            <Input placeholder="Access Key Id" type="password" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            label="Region"
            name={[...baseName, 'region']}
            rules={[{required: true, message: 'Region is required'}]}
          >
            <Input placeholder="Region" />
          </Form.Item>
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Secret Access Key"
            name={[...baseName, 'secretAccessKey']}
            rules={[{required: true, message: 'Secret Access Key is required'}]}
          >
            <Input placeholder="Secret Access Key" type="password" />
          </Form.Item>
        </Col>
      </Row>
    </>
  );
};

export default AwsXRay;
