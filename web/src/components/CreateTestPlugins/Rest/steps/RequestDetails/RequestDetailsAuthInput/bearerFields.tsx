import {Form, Input} from 'antd';
import * as S from '../RequestDetails.styled';

export const bearerFields: React.ReactElement = (
  <S.Row style={{width: '100%'}}>
    <Form.Item data-cy="bearer-token" name={['auth', 'bearer', 'token']} label="Token" rules={[{required: true}]}>
      <Input />
    </Form.Item>
  </S.Row>
);
