import {Checkbox, Col, Form, Input, Row} from 'antd';

import {FileUpload} from 'components/Inputs';
import {SupportedDataStoresToDefaultEndpoint, SupportedDataStoresToName} from 'constants/DataStore.constants';
import {SupportedDataStores, TDraftDataStore} from 'types/DataStore.types';
import * as S from '../../DataStorePluginForm.styled';
import AddressesList from './AddressesList';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const OpenSearch = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = Form.useWatch('dataStoreType', form) || SupportedDataStores.OpenSearch;
  const baseName = ['dataStore', dataStoreType];
  const endpointPlaceholder = SupportedDataStoresToDefaultEndpoint[dataStoreType];

  return (
    <>
      <S.Title>Provide the connection info for {SupportedDataStoresToName[dataStoreType]}</S.Title>
      <DataStoreDocsBanner dataStoreType={dataStoreType} />
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Index"
            name={[...baseName, 'index']}
            rules={[{required: true, message: 'Index is required'}]}
          >
            <Input placeholder="Index" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <S.ItemListLabel>Addresses</S.ItemListLabel>
          <Form.List name={[...baseName, 'addresses']}>
            {(fields, {add, remove}) => (
              <AddressesList fields={fields} add={add} remove={remove} placeholder={endpointPlaceholder} />
            )}
          </Form.List>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item label="Username" name={[...baseName, 'username']}>
            <Input placeholder="Username" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item label="Password" name={[...baseName, 'password']}>
            <Input placeholder="Password" type="password" />
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item name={[...baseName, 'insecureSkipVerify']} valuePropName="checked">
            <Checkbox>Enable TLS but not verify the certificate</Checkbox>
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item label="Upload CA file" name={[...baseName, 'certificateFile']}>
            <FileUpload accept="" />
          </Form.Item>
        </Col>
      </Row>
    </>
  );
};

export default OpenSearch;
