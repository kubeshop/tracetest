import {Col, Row} from 'antd';
import Ingestor from './Ingestor';

const Agent = () => (
  <Row gutter={[16, 16]}>
    <Col span={24}>
      <Ingestor />
    </Col>
  </Row>
);

export default Agent;
