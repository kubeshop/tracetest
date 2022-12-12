import {Checkbox, Col, Form, Input, Row, Select, Space, Switch} from 'antd';
import {useState} from 'react';

import RequestDetailsAuthInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import {TDraftDataStore} from 'types/Config.types';
import * as S from './GrcpClient.styled';
import GrpcClientSecure from './GrpcClientSecure';

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
  const [isSecure, setIsSecure] = useState(false);
  const dataStoreType = form.getFieldValue('dataStoreType');
  const baseName = ['dataStore', dataStoreType];

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="Endpoint"
            name={[...baseName, 'endpoint']}
            rules={[{required: true, message: 'Endpoint is required'}]}
          >
            <Input placeholder="Enter an endpoint" />
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
            <Input placeholder="Enter a read buffer size" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item label="Write Buffer Size" name={[...baseName, 'writeBufferSize']}>
            <Input placeholder="Enter a write buffer size" />
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
          <RequestDetailsAuthInput hasBaseApikeyFields name={[...baseName, 'auth']} />
        </Col>
        <Col span={12}>
          <RequestDetailsHeadersInput initialValue={HEADER_DEFAULT_VALUES} name={[...baseName, 'headers']} />
        </Col>
      </Row>

      <Space>
        <Switch
          onChange={checked => {
            form.setFieldsValue({dataStore: {[dataStoreType]: {tls: {insecure: !checked}}}});
            setIsSecure(checked);
          }}
        />{' '}
        Secure options
        <Form.Item hidden initialValue name={[...baseName, 'tls', 'insecure']} valuePropName="checked">
          <Checkbox>Insecure</Checkbox>
        </Form.Item>
      </Space>

      {isSecure && <GrpcClientSecure baseName={baseName} />}
    </>
  );
};

export default GrpcClient;
