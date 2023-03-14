import {Col, Form, Input, Row} from 'antd';
import {SupportedDemosFormField} from 'types/Settings.types';

const OtelFields = () => {
  const baseName = [SupportedDemosFormField.OpentelemetryStore, SupportedDemosFormField.OpentelemetryStore];

  return (
    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="Frontend Endpoint" name={[...baseName, 'frontendEndpoint']}>
          <Input placeholder="http://otel-frontend.otel-demo:8084" />
        </Form.Item>

        <Form.Item label="Cart Endpoint" name={[...baseName, 'cartEndpoint']}>
          <Input placeholder="http://otel-cartservice.otel-demo:7070" />
        </Form.Item>
      </Col>

      <Col span={12}>
        <Form.Item label="Product Catalog Endpoint" name={[...baseName, 'productCatalogEndpoint']}>
          <Input placeholder="http://otel-productcatalogservice.otel-demo:3550" />
        </Form.Item>

        <Form.Item label="Checkout Endpoint" name={[...baseName, 'checkoutEndpoint']}>
          <Input placeholder="http://otel-checkoutservice.otel-demo:5050" />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default OtelFields;
