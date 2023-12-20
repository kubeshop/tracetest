import {Col} from 'antd';
import UrlCodeBlock from 'components/CodeBlock/UrlCodeBlock';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import {INGESTOR_ENDPOINT_URL} from 'constants/Common.constants';
import * as S from './OpenTelemetryCollector.styled';

const Ingestor = () => (
  <S.Container>
    <S.Title>Ingestor Endpoint</S.Title>
    <S.Description>
      Tracetest exposes trace ingestion endpoints on ports 4317 for gRPC and 4318 for HTTP. Use the Tracetest Server’s
      hostname and port to connect. For example, with Docker use tracetest:4317 for gRPC.
    </S.Description>

    <Col span={16}>
      <S.UrlEntry>
        gRPC <UrlCodeBlock value="tracetest:4317" minHeight="35px" maxHeight="35px" language="bash" />
      </S.UrlEntry>
    </Col>
    <Col span={16}>
      <S.UrlEntry>
        HTTP <UrlCodeBlock value="tracetest:4318/v1/traces" minHeight="35px" maxHeight="35px" language="bash" />
      </S.UrlEntry>
    </Col>
    <Col span={16}>
      <DocsBanner>
        Need more information about setting up ingestion endpoint?{' '}
        <a target="_blank" href={INGESTOR_ENDPOINT_URL}>
          Go to our docs
        </a>
      </DocsBanner>
    </Col>
  </S.Container>
);

export default Ingestor;
