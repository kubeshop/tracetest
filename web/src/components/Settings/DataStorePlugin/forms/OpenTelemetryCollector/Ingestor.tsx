import {Form, Switch} from 'antd';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import {INGESTOR_ENDPOINT_URL} from 'constants/Common.constants';
import {TCollectorDataStores, TDraftDataStore} from 'types/DataStore.types';
import * as S from './OpenTelemetryCollector.styled';

const Ingestor = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = Form.useWatch('dataStoreType', form) as TCollectorDataStores;
  const baseName = ['dataStore', dataStoreType];

  return (
    <S.Container>
      <S.Description>
        Tracetest easily integrates with any distributed tracing solution via the OpenTelemetry Collector. It allows
        your current tracing system to send only Tracetest spans while the rest go to your chosen backend.
      </S.Description>
      <S.Title>Ingestor Endpoint</S.Title>
      <S.Description>
        Tracetest exposes trace ingestion endpoints on ports 4317 for gRPC and 4318 for HTTP. Turn on the Tracetest
        ingestion endpoint to start sending traces. Use the Tracetest Server’s hostname and port to connect.For example,
        with Docker use tracetest:4317 for gRPC.
      </S.Description>
      <S.SwitchContainer>
        <Form.Item name={[...baseName, 'isIngestorEnabled']} valuePropName="checked">
          <Switch />
        </Form.Item>
        Enable Tracetest ingestion endpoint
      </S.SwitchContainer>
      <DocsBanner>
        Need more information about setting up ingestion endpoint?{' '}
        <a target="_blank" href={INGESTOR_ENDPOINT_URL}>
          Go to our docs
        </a>
      </DocsBanner>
    </S.Container>
  );
};

export default Ingestor;
