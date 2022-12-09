import {Checkbox, Col, Form, Input, Row} from 'antd';

interface IProps {
  baseName: string[];
}

const GrpcClientSecure = ({baseName}: IProps) => (
  <>
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <Form.Item name={[...baseName, 'tls', 'insecureSkipVerify']} valuePropName="checked">
          <Checkbox>Enable TLS but not verify the certificate</Checkbox>
        </Form.Item>
      </Col>
    </Row>
    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="Server Name" name={[...baseName, 'tls', 'serverName']}>
          <Input placeholder="Enter a server name" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="CA file" name={[...baseName, 'tls', 'settings', 'cAFile']}>
          <Input placeholder="Enter a CA file" />
        </Form.Item>
      </Col>
    </Row>

    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="Cert file" name={[...baseName, 'tls', 'settings', 'certFile']}>
          <Input placeholder="Enter a Cert file" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="Key file" name={[...baseName, 'tls', 'settings', 'keyFile']}>
          <Input placeholder="Enter a Key file" />
        </Form.Item>
      </Col>
    </Row>

    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="TLS Min version" name={[...baseName, 'tls', 'settings', 'minVersion']}>
          <Input placeholder="Enter a min version" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="TLS Max version" name={[...baseName, 'tls', 'settings', 'maxVersion']}>
          <Input placeholder="Enter a max version" />
        </Form.Item>
      </Col>
    </Row>
  </>
);

export default GrpcClientSecure;
