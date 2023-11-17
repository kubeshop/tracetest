import {Col, Row} from 'antd';
import Ingestor from './Ingestor';
import Configuration from '../OpenTelemetryCollector/Configuration';

const Agent = () => (
  <>
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <Ingestor />
      </Col>
    </Row>
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <Configuration />
      </Col>
    </Row>
  </>
);

export default Agent;
