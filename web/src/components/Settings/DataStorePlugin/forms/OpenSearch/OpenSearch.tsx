import {Form, Input} from 'antd';
import {SupportedDataStores} from 'types/Config.types';
import * as S from '../../DataStorePluginForm.styled';
import AddressesList from './AddressesList';

const OpenSearch = () => {
  const baseName = ['dataStore', SupportedDataStores.OpenSearch];

  return (
    <S.FormContainer>
      <S.FormColumn>
        <Form.Item
          label="Username"
          name={[...baseName, 'username']}
          rules={[{required: true, message: 'Username is required'}]}
        >
          <Input placeholder="Username" />
        </Form.Item>
        <Form.Item
          label="Password"
          name={[...baseName, 'password']}
          rules={[{required: true, message: 'Password is required'}]}
        >
          <Input placeholder="Password" type="password" />
        </Form.Item>
      </S.FormColumn>
      <S.FormColumn>
        <Form.Item label="Index" name={[...baseName, 'index']} rules={[{required: true, message: 'Index is required'}]}>
          <Input placeholder="Index" />
        </Form.Item>
        <div>
          <S.ItemListLabel>Addresses</S.ItemListLabel>
          <Form.List name={[...baseName, 'addresses']}>
            {(fields, {add, remove}) => <AddressesList fields={fields} add={add} remove={remove} />}
          </Form.List>
        </div>
      </S.FormColumn>
    </S.FormContainer>
  );
};

export default OpenSearch;
