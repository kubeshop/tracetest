import {Col, Form, Input, Row} from 'antd';

const PokeshopFields = () => (
  <Row gutter={[16, 16]}>
    <Col span={12}>
      <Form.Item label="HTTP Endpoint" name="pokeshopHttp">
        <Input placeholder="http://demo-pokemon-api.demo" />
      </Form.Item>
    </Col>

    <Col span={12}>
      <Form.Item label="GRPC Endpoint" name="pokeshopGrpc">
        <Input placeholder="demo-pokemon-api.demo:8082" />
      </Form.Item>
    </Col>
  </Row>
);

export default PokeshopFields;
