import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStores} from 'types/DataStore.types';

const AzureAppInsights = () => {
  const baseName = ['dataStore', SupportedDataStores.AzureAppInsights];

  return (
    <Row gutter={[16, 16]}>
      <Col span={16}>
        <Form.Item
          label="Application Insights Artifact Resource Name Id"
          name={[...baseName, 'resourceArmId']}
          rules={[{required: true, message: 'Access Key Id is required'}]}
        >
          <Input placeholder="/subscriptions/<sid>/resourceGroups/<rg>/providers/<providerName>/<resourceType>/<resourceName>" />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default AzureAppInsights;
