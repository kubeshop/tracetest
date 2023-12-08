import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import {SupportedDataStores} from 'types/DataStore.types';
import * as S from '../../DataStorePluginForm.styled';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const SumoLogic = () => {
  const baseName = ['dataStore', SupportedDataStores.SumoLogic];

  return (
    <>
      <S.Title>Provide the connection info for {SupportedDataStoresToName[SupportedDataStores.SumoLogic]}</S.Title>
      <DataStoreDocsBanner dataStoreType={SupportedDataStores.SumoLogic} />

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item label="URL" name={[...baseName, 'url']} rules={[{required: true, message: 'URL is required'}]}>
            <Input placeholder="Enter a URL" />
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
