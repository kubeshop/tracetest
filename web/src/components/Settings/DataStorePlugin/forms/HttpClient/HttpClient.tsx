import {Checkbox, Col, Form, Input, Row, Space, Switch} from 'antd';
import RequestDetailsHeadersInput from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import {TDraftDataStore} from 'types/Config.types';
import GrpcClientSecure from '../GrpcClient/GrpcClientSecure';

const HEADER_DEFAULT_VALUES = [{key: '', value: ''}];

const HttpClient = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = form.getFieldValue('dataStoreType');
  const baseName = ['dataStore', dataStoreType, 'http'];
  const insecureValue = Form.useWatch([...baseName, 'tls', 'insecure'], form) ?? true;

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          <Form.Item
            label="URL"
            name={[...baseName, 'url']}
            rules={[{required: true, message: 'URL is required'}]}
          >
            <Input placeholder="Enter a URL" />
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={12}>
          <RequestDetailsHeadersInput initialValue={HEADER_DEFAULT_VALUES} name={[...baseName, 'rawHeaders']} />
        </Col>
      </Row>

      <Space>
        <Switch
          onChange={checked => {
            form.setFieldsValue({
              dataStore: {
                name: dataStoreType,
                type: dataStoreType,
                [dataStoreType]: {http: {tls: {insecure: !checked}}},
              },
            });
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

export default HttpClient;
