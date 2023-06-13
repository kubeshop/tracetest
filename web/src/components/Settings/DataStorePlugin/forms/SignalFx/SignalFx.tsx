import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStores} from 'types/DataStore.types';
import {SupportedDataStoresToName} from 'constants/DataStore.constants';
import * as S from '../../DataStorePluginForm.styled';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const SignalFx = () => {
  const baseName = ['dataStore', SupportedDataStores.SignalFX];

  return (
    <>
      <S.Title>Provide the connection info for {SupportedDataStoresToName[SupportedDataStores.SignalFX]}</S.Title>
      <DataStoreDocsBanner dataStoreType={SupportedDataStores.SignalFX} />
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Realm"
            name={[...baseName, 'realm']}
            rules={[{required: true, message: 'Realm is required'}]}
          >
            <Input placeholder="us1" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item
            label="Token"
            name={[...baseName, 'token']}
            rules={[{required: true, message: 'Token is required'}]}
          >
            <Input placeholder="Your token" type="password" />
          </Form.Item>
        </Col>
      </Row>
    </>
  );
};

export default SignalFx;
