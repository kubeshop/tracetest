import {Checkbox, Col, Form, Input, Row, Select, Space, Switch} from 'antd';
import {useCallback} from 'react';
import {SupportedDataStoresToDefaultEndpoint, SupportedDataStoresToName} from 'constants/DataStore.constants';
import {useDataStore} from 'providers/DataStore/DataStore.provider';
import DataStoreService from 'services/DataStore.service';
import {SupportedDataStores, TDraftDataStore} from 'types/DataStore.types';
import * as FS from '../../DataStorePluginForm.styled';
import * as S from './GrcpClient.styled';
import GrpcClientSecure from './GrpcClientSecure';
import DataStoreDocsBanner from '../../../DataStoreDocsBanner/DataStoreDocsBanner';

const COMPRESSION_LIST = [
  {name: 'none', value: 'none'},
  {name: 'gzip', value: 'gzip'},
  {name: 'zlib', value: 'zlib'},
  {name: 'deflate', value: 'deflate'},
  {name: 'snappy', value: 'snappy'},
  {name: 'zstd', value: 'zstd'},
] as const;

const HEADER_DEFAULT_VALUES = [{key: '', value: ''}];

const GrpcClient = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = form.getFieldValue('dataStoreType') as SupportedDataStores;
  const baseName = ['dataStore', dataStoreType, 'grpc'];
  const insecureValue = Form.useWatch([...baseName, 'tls', 'insecure'], form) ?? true;
  const {onIsFormValid} = useDataStore();

  const onValidation = useCallback(
    async (draft: TDraftDataStore) => {
      const isValid = await DataStoreService.validateDraft(draft);
      onIsFormValid(isValid);
    },
    [onIsFormValid]
  );

  return (
    <>
      <FS.Title>Provide the connection info for {SupportedDataStoresToName[dataStoreType]}</FS.Title>
      <DataStoreDocsBanner dataStoreType={dataStoreType} />
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Endpoint"
            name={[...baseName, 'endpoint']}
            rules={[{required: true, message: 'Endpoint is required'}]}
          >
            <Input placeholder={SupportedDataStoresToDefaultEndpoint[dataStoreType]} />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item label="Compression" name={[...baseName, 'compression']}>
            <Select placeholder="None" allowClear onClear={() => form.resetFields(['compression'])}>
              {COMPRESSION_LIST.map(({name, value}) => (
                <Select.Option key={value} value={value}>
                  {name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item label="Read Buffer Size" name={[...baseName, 'readBufferSize']}>
            <Input placeholder="Enter a read buffer size" type="number" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item label="Write Buffer Size" name={[...baseName, 'writeBufferSize']}>
            <Input placeholder="Enter a write buffer size" type="number" />
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item label="Balancer Name" name={[...baseName, 'balancerName']}>
            <Input placeholder="Enter a balancer name" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <S.ChecksContainer>
            <Form.Item name={[...baseName, 'waitForReady']} valuePropName="checked">
              <Checkbox>Wait For Ready State</Checkbox>
            </Form.Item>
          </S.ChecksContainer>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <FS.HeadersInput initialValue={HEADER_DEFAULT_VALUES} name={[...baseName, 'rawHeaders']} />
        </Col>
      </Row>

      <Space>
        <Switch
          onChange={checked => {
            form.setFieldsValue({
              dataStore: {
                name: dataStoreType,
                type: dataStoreType,
                [dataStoreType]: {grpc: {tls: {insecure: !checked}}},
              },
            });
            onValidation(form.getFieldsValue());
          }}
          checked={!insecureValue}
        />{' '}
        Secure options
        <Form.Item hidden initialValue name={[...baseName, 'tls', 'insecure']} valuePropName="checked">
          <Checkbox>Insecure</Checkbox>
        </Form.Item>
      </Space>

      {!insecureValue && <GrpcClientSecure baseName={baseName} />}
    </>
  );
};

export default GrpcClient;
