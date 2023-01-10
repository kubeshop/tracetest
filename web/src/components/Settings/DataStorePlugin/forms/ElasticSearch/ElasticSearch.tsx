import {Col, Form, Input, Row} from 'antd';
import {SupportedDataStores} from 'types/Config.types';
import RequestDetailsFileInput from '../../../../CreateTestPlugins/Grpc/steps/RequestDetails/RequestDetailsFileInput';
import * as S from '../../DataStorePluginForm.styled';
import AddressesList from './AddressesList';

const OpenSearch = () => {
  const baseName = ['dataStore', SupportedDataStores.OpenSearch];

  return (
    <>
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
            {(fields, {add, remove}) => <AddressesList fields={fields} add={add} remove={remove} />}
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
          <Form.Item label="Upload CA file" name="certificateFile" required={false}>
            <RequestDetailsFileInput accept="" />
          </Form.Item>
        </Col>
      </Row>
    </>
  );
};

export default OpenSearch;
