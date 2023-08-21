import {Col, Form, Input, Row} from 'antd';
import {SupportedDemos} from 'types/Settings.types';

const PokeshopFields = () => {
  const baseName = [SupportedDemos.Pokeshop, SupportedDemos.Pokeshop];

  return (
    <Row gutter={[16, 16]}>
      <Col span={12}>
        <Form.Item label="HTTP Endpoint" name={[...baseName, 'httpEndpoint']}>
          <Input placeholder="http://demo-pokemon-api.demo" />
        </Form.Item>
      </Col>

      <Col span={12}>
        <Form.Item label="GRPC Endpoint" name={[...baseName, 'grpcEndpoint']}>
          <Input placeholder="demo-pokemon-api.demo:8082" />
        </Form.Item>
      </Col>

      <Col span={12}>
        <Form.Item label="Kafka Broker" name={[...baseName, 'kafkaBroker']}>
          <Input placeholder="stream.demo:9092" />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default PokeshopFields;
