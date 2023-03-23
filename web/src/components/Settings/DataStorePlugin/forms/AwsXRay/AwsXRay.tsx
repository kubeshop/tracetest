import {Checkbox, Col, Form, Input, Row} from 'antd';
import {SupportedDataStores, TDraftDataStore} from 'types/DataStore.types';

const AwsXRay = () => {
  const baseName = ['dataStore', SupportedDataStores.AWSXRay];
  const form = Form.useFormInstance<TDraftDataStore>();
  const useDefaultAuth = Form.useWatch([...baseName, 'useDefaultAuth'], form) ?? false;

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item name={[...baseName, 'useDefaultAuth']} valuePropName="checked">
            <Checkbox>Use Default AWS Authentication</Checkbox>
          </Form.Item>
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Access Key Id"
            name={[...baseName, 'accessKeyId']}
            rules={[{required: !useDefaultAuth, message: 'Access Key Id is required'}]}
          >
            <Input placeholder="Access Key Id" type="password" disabled={useDefaultAuth} />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            label="Secret Access Key"
            name={[...baseName, 'secretAccessKey']}
            rules={[{required: !useDefaultAuth, message: 'Secret Access Key is required'}]}
          >
            <Input placeholder="Secret Access Key" type="password" disabled={useDefaultAuth} />
          </Form.Item>
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item label="Session Token" name={[...baseName, 'sessionToken']}>
            <Input placeholder="Session Token" type="password" disabled={useDefaultAuth} />
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
    </>
  );
};

export default AwsXRay;
