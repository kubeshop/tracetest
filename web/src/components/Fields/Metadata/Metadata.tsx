import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor} from 'components/Inputs';
import * as S from './Metadata.styled';

const Metadata = () => (
  <Form.Item label="Metadata">
    <Form.List name="metadata">
      {(fields, {add, remove}) => (
        <>
          {fields.map(field => (
            <S.HeaderContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder="Key" />
              </Form.Item>

              <Form.Item name={[field.name, 'value']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder="Value" />
              </Form.Item>

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
