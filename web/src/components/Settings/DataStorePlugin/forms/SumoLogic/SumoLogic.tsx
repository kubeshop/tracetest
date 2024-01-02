import {Col, Form, Input, Row, Typography} from 'antd';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import * as S from '../../DataStorePluginForm.styled';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const SUMO_LOGIC_ENDPOINTS_DOCS_URL =
  'https://help.sumologic.com/docs/api/getting-started/#sumo-logic-endpoints-by-deployment-and-firewall-security';

const SumoLogic = () => {
  const baseName = ['dataStore', SupportedDataStores.SumoLogic];

  return (
    <>
      <S.Title>Provide the connection info for {SupportedDataStoresToName[SupportedDataStores.SumoLogic]}</S.Title>
      <DataStoreDocsBanner dataStoreType={SupportedDataStores.SumoLogic} />

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Sumo Logic API endpoint"
            name={[...baseName, 'url']}
            rules={[{required: true, message: 'URL is required'}]}
            help={
              <>
                <Typography.Text type="secondary">
                  This URL will be different based on your location. Use the Sumo Logic API endpoint as documented on
                </Typography.Text>
                &nbsp;
                <Typography.Link href={SUMO_LOGIC_ENDPOINTS_DOCS_URL} target="_blank">
                  &quot;Which endpoint should I use&quot;
                </Typography.Link>
              </>
            }
          >
            <Input placeholder="https://api.sumologic.com/api/" />
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Access ID"
            name={[...baseName, 'accessID']}
            rules={[{required: true, message: 'Access ID is required'}]}
          >
            <Input placeholder="Access ID" type="password" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            label="Access Key"
            name={[...baseName, 'accessKey']}
            rules={[{required: true, message: 'Access Key is required'}]}
          >
            <Input placeholder="Access Key" type="password" />
          </Form.Item>
        </Col>
      </Row>
    </>
  );
};

export default SumoLogic;
