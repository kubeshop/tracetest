import {Form, Input} from 'antd';
import {IPostmanValues} from 'types/Test.types';
import {CollectionFileField} from './fields/CollectionFileField';
import {EnvFileField} from './fields/EnvFileField';
import {SelectTestFromCollection} from './fields/SelectTestFromCollection';
import * as S from './Postman.styled';

const Postman = () => {
  const form = Form.useFormInstance<IPostmanValues>();

  return (
    <S.FieldsContainer>
      <Form.Item name="requests" hidden>
        <Input type="hidden" />
      </Form.Item>
      <Form.Item name="variables" hidden>
        <Input type="hidden" />
      </Form.Item>
      <CollectionFileField form={form} />
      <EnvFileField form={form} />
      <SelectTestFromCollection form={form} />
    </S.FieldsContainer>
  );
};

export default Postman;
