import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import * as S from './Metadata.styled';
import SingleLine from '../../Inputs/SingleLine';

const Metadata = () => (
  <Form.Item>
    <Form.List name="metadata" initialValue={[{key: '', value: ''}]}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.HeaderContainer key={field.name}>
              <S.Item>
                <Form.Item name={[field.name, 'key']} noStyle>
                  <SingleLine placeholder={`Key ${index + 1}`} />
                </Form.Item>
              </S.Item>

              <S.Item>
                <Form.Item name={[field.name, 'value']} noStyle>
                  <SingleLine placeholder={`Value ${index + 1}`} />
                </Form.Item>
              </S.Item>

              <Form.Item noStyle>
                <Button
                  icon={<S.DeleteIcon />}
                  onClick={() => remove(field.name)}
                  style={{marginLeft: 12}}
                  type="link"
                />
              </Form.Item>
            </S.HeaderContainer>
          ))}

          <Button
            data-cy="add-metadata"
            icon={<PlusOutlined />}
            onClick={() => add()}
            style={{fontWeight: 600, height: 'auto', padding: 0}}
            type="link"
          >
            Add Entry
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default Metadata;
