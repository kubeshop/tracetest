import {Form, Switch} from 'antd';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import {TCollectorDataStores, TDraftDataStore} from 'types/DataStore.types';
import {INGESTOR_ENDPOINT_URL} from 'constants/Common.constants';
import * as S from './Agent.styled';

const Ingestor = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = Form.useWatch('dataStoreType', form) as TCollectorDataStores;
  const baseName = ['dataStore', dataStoreType];

  return (
    <S.Container>
      <S.Description>
        The Tracetest Agent can be used to collect OpenTelemetry trace information from the host it is running on.
      </S.Description>
      <S.Title>Ingestor Endpoint</S.Title>
      <S.Description>
        The Tracetest Agent exposes trace ingestion endpoints on ports 4317 for gRPC and 4318 for HTTP. Turn on the
        Tracetest ingestion endpoint to start sending traces. Use your local hostname and port to connect. For example,
        localhost:4317 for gRPC.
      </S.Description>
      <S.SwitchContainer>
        <Form.Item name={[...baseName, 'isIngestorEnabled']} valuePropName="checked">
          <Switch />
        </Form.Item>
        <label htmlFor={`data-store_dataStore_${dataStoreType}_isIngestorEnabled`}>
          Enable Tracetest ingestion endpoint
        </label>
      </S.SwitchContainer>
      <DocsBanner>
        Need more information about setting up the agent ingestion endpoint?{' '}
        <a target="_blank" href={INGESTOR_ENDPOINT_URL}>
          Go to our docs
        </a>
      </DocsBanner>
    </S.Container>
  );
};

export default Ingestor;
