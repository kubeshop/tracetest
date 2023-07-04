import {Col, Form, Radio, Row} from 'antd';
import {SupportedClientTypes, SupportedDataStores, TDraftDataStore} from 'types/DataStore.types';
import GrpcClient from '../GrpcClient';
import HttpClient from '../HttpClient';

const FieldsFormMap = {
  [SupportedClientTypes.GRPC]: GrpcClient,
  [SupportedClientTypes.HTTP]: HttpClient,
} as const;

const BaseClient = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = form.getFieldValue('dataStoreType') as SupportedDataStores;
  const baseName = ['dataStore', dataStoreType];
  const type = (Form.useWatch([...baseName, 'type'], form) || SupportedClientTypes.GRPC) as SupportedClientTypes;
  const Component = FieldsFormMap[type];

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={12}>
          Connection type:
          <Form.Item name={[...baseName, 'type']}>
            <Radio.Group defaultValue={type}>
              <Radio value={SupportedClientTypes.GRPC}>gRPC</Radio>
              <Radio value={SupportedClientTypes.HTTP}>Http</Radio>
            </Radio.Group>
          </Form.Item>
        </Col>
      </Row>
      <Component />
    </>
  );
};

export default BaseClient;
