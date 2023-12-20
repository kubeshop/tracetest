import {Col, Row} from 'antd';
import {withCustomization} from 'providers/Customization';
import Ingestor from './Ingestor';
import Configuration from './Configuration';

const OpenTelemetryCollector = () => (
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

export default withCustomization(OpenTelemetryCollector, 'openTelemetryCollector');
