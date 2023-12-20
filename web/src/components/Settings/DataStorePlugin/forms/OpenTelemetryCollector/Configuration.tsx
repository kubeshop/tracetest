import {Col, Collapse, Form, Row, Typography} from 'antd';
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
    <S.CollapseContainer>
      <Collapse ghost>
        <Collapse.Panel
          header={<Typography.Title level={3}>OpenTelemetry Collector Configuration</Typography.Title>}
          key="1"
        >
          <S.SubtitleContainer>
            <S.Description>
              The OpenTelemetry Collector configuration below is a sample. Your config file layout should look the same.
              Make sure to use your own API keys or tokens as explained. Copy the config sample below, paste it into
              your own OpenTelemetry Collector config and apply it.
            </S.Description>
          </S.SubtitleContainer>

          <Row>
            <Col span={16}>
              <S.CodeContainer data-cy="file-viewer-code-container">
                <FramedCodeBlock value={example} language="yaml" title="Collector Configuration" />
              </S.CodeContainer>
            </Col>
          </Row>
        </Collapse.Panel>
      </Collapse>

      <DataStoreDocsBanner dataStoreType={dataStoreType} />
    </S.CollapseContainer>
  );
};

export default Configuration;
