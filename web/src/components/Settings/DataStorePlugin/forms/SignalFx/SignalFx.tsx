import {Form, Input} from 'antd';
import {SupportedDataStores} from 'types/Config.types';
import * as S from '../../DataStorePluginForm.styled';

const SignalFx = () => {
  const baseName = ['dataStore', SupportedDataStores.SignalFX];

  return (
    <S.FormContainer>
      <S.FormColumn>
        <Form.Item label="Realm" name={[...baseName, 'realm']} rules={[{required: true, message: 'Realm is required'}]}>
          <Input placeholder="Realm" />
        </Form.Item>
      </S.FormColumn>
      <S.FormColumn>
        <Form.Item label="Token" name={[...baseName, 'token']} rules={[{required: true, message: 'Token is required'}]}>
          <Input placeholder="Token" type="password" />
        </Form.Item>
      </S.FormColumn>
    </S.FormContainer>
  );
};

export default SignalFx;
