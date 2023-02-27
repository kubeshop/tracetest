import {Col, Form, Input, Row} from 'antd';

const OtelFields = () => (
  <Row gutter={[16, 16]}>
    <Col span={12}>
      <Form.Item label="Frontend Endpoint" name={['demo', 'otelFrontend']}>
        <Input placeholder="Enter an endpoint" />
      </Form.Item>

      <Form.Item label="Cart Endpoint" name={['demo', 'otelCart']}>
        <Input placeholder="Enter an endpoint" />
      </Form.Item>
    </Col>

    <Col span={12}>
      <Form.Item label="Product Catalog Endpoint" name={['demo', 'otelProductCatalog']}>
        <Input placeholder="Enter an endpoint" />
      </Form.Item>

      <Form.Item label="Checkout Endpoint" name={['demo', 'otelCheckout']}>
        <Input placeholder="Enter an endpoint" />
      </Form.Item>
    </Col>
  </Row>
);

export default OtelFields;
