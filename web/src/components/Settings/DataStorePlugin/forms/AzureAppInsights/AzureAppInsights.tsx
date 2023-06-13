import {Checkbox, Col, Form, Input, Radio, Row} from 'antd';
import {ConnectionTypes, SupportedDataStores, TDraftDataStore} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import OpenTelemetryCollector from '../OpenTelemetryCollector/OpenTelemetryCollector';
import * as S from '../../DataStorePluginForm.styled';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const AzureAppInsights = () => {
  const baseName = ['dataStore', SupportedDataStores.AzureAppInsights];
  const form = Form.useFormInstance<TDraftDataStore>();
  const connectionType = Form.useWatch([...baseName, 'connectionType'], form);
  const useAzureActiveDirectoryAuth = Form.useWatch([...baseName, 'useAzureActiveDirectoryAuth'], form);

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          Connection type:
          <Form.Item name={[...baseName, 'connectionType']}>
            <Radio.Group>
              <Radio value={ConnectionTypes.Direct}>Direct Connection</Radio>
              <Radio value={ConnectionTypes.Collector}>Open Telemetry Collector</Radio>
            </Radio.Group>
          </Form.Item>
        </Col>
      </Row>
      {(connectionType === ConnectionTypes.Collector && <OpenTelemetryCollector />) || (
        <>
          <S.Title>
            Provide the connection info for {SupportedDataStoresToName[SupportedDataStores.AzureAppInsights]}
          </S.Title>
          <DataStoreDocsBanner dataStoreType={SupportedDataStores.AzureAppInsights} />
          <Row gutter={[16, 16]}>
            <Col span={16}>
              <Form.Item
                label="Application Insights Artifact Resource Name Id"
                name={[...baseName, 'resourceArmId']}
                rules={[{required: true, message: 'ARM Id is required'}]}
              >
                <Input placeholder="/subscriptions/<sid>/resourceGroups/<rg>/providers/<providerName>/<resourceType>/<resourceName>" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={12}>
              <Form.Item name={[...baseName, 'useAzureActiveDirectoryAuth']} valuePropName="checked">
                <Checkbox>Use Azure Active Directory Auth</Checkbox>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={[16, 16]}>
            <Col span={16}>
              <Form.Item
                label="Access Token"
                name={[...baseName, 'accessToken']}
                rules={[{required: !useAzureActiveDirectoryAuth, message: 'Access Token is required'}]}
              >
                <Input disabled={useAzureActiveDirectoryAuth} type="password" placeholder="your access token" />
              </Form.Item>
            </Col>
          </Row>
        </>
      )}
    </>
  );
};

export default AzureAppInsights;
