import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './RequestDetails.styled';

interface IProps {
  name?: string[];
  className?: string;
}
const RequestDetailsHeadersInput = ({
  name = ['headers'],
  className = '',
}: IProps) => (
  <Form.Item className={`input-headers ${className}`} label="Message Headers" shouldUpdate>
    <Form.List name={name.length === 1 ? name[0] : name}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.HeaderContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder={`Header Key ${index + 1}`} />
              </Form.Item>

              <Form.Item name={[field.name, 'value']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder={`Header Value ${index + 1}`} />
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
            data-cy="add-header"
            icon={<PlusOutlined />}
            onClick={() => add()}
            style={{fontWeight: 600, height: 'auto', padding: 0}}
            type="link"
          >
            Add Header
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default RequestDetailsHeadersInput;
