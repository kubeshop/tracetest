import {Form} from 'antd';
import {CollectorConfigMap} from 'constants/CollectorConfig.constants';
import {TCollectorDataStores, TDraftDataStore} from 'types/DataStore.types';
import {FramedCodeBlock} from 'components/CodeBlock';
import * as S from './OpenTelemetryCollector.styled';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const Configuration = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = Form.useWatch('dataStoreType', form) as TCollectorDataStores;
  const example = CollectorConfigMap[dataStoreType!];

  return (
    <>
      <S.SubtitleContainer>
        <S.Title>OpenTelemetry Collector Configuration</S.Title>
        <S.Description>
          The OpenTelemetry Collector configuration below is a sample. Your config file layout should look the same.
          Make sure to use your own API keys or tokens as explained. Copy the config sample below, paste it into your
          own OpenTelemetry Collector config and apply it.
        </S.Description>
      </S.SubtitleContainer>
      <S.CodeContainer data-cy="file-viewer-code-container">
        <FramedCodeBlock value={example} language="yaml" title="Collector Configuration" />
      </S.CodeContainer>
      <DataStoreDocsBanner dataStoreType={dataStoreType} />
    </>
  );
};

export default Configuration;
