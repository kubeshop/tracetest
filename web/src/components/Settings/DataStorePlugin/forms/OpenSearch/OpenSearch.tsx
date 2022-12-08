import {Form, Input} from 'antd';
import * as S from '../../DataStorePluginForm.styled';
import AddressesList from './AddressesList';

const OpenSearch = () => {
  return (
    <S.FormContainer>
      <S.FormColumn>
        <Form.Item label="Username" name="username">
          <Input placeholder="Username" />
        </Form.Item>
        <Form.Item label="Password" name="password">
          <Input placeholder="Password" type="password" />
        </Form.Item>
      </S.FormColumn>
      <S.FormColumn>
        <Form.Item label="Index" name="index">
          <Input placeholder="Index" />
        </Form.Item>
        <div>
          <S.ItemListLabel>Addresses</S.ItemListLabel>
          <Form.List name="addresses">
            {(fields, {add, remove}) => <AddressesList fields={fields} add={add} remove={remove} />}
          </Form.List>
        </div>
      </S.FormColumn>
    </S.FormContainer>
  );
};

export default OpenSearch;
