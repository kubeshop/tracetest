import {Checkbox, Col, Form, Input, Row} from 'antd';
import RequestDetailsFileInput from 'components/CreateTestPlugins/Grpc/steps/RequestDetails/RequestDetailsFileInput';

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
        <Form.Item label="Upload CA file" name={[...baseName, 'fileCA']}>
          <RequestDetailsFileInput accept="" />
        </Form.Item>
      </Col>
    </Row>

    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="Upload Cert file" name={[...baseName, 'fileCert']}>
          <RequestDetailsFileInput accept="" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="Upload Key file" name={[...baseName, 'fileKey']}>
          <RequestDetailsFileInput accept="" />
        </Form.Item>
      </Col>
    </Row>

    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="TLS Min version" name={[...baseName, 'tls', 'minVersion']}>
          <Input placeholder="Enter a min version" />
        </Form.Item>
      </Col>
      <Col span={12}>
        <Form.Item label="TLS Max version" name={[...baseName, 'tls', 'maxVersion']}>
          <Input placeholder="Enter a max version" />
        </Form.Item>
      </Col>
    </Row>
  </>
);

export default GrpcClientSecure;
