import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import {DEFAULT_HEADERS, IKeyValue} from 'constants/Test.constants';
import Editor from 'components/Editor';
import {SupportedEditors} from 'constants/Editor.constants';
import * as S from './RequestDetails.styled';

interface IProps {
  initialValue?: IKeyValue[];
  name?: string[];
  unit?: string;
  label?: string;
  className?: string;
}
const RequestDetailsHeadersInput = ({
  unit = 'Header',
  name = ['headers'],
  initialValue = DEFAULT_HEADERS,
  label = 'Header',
  className = '',
}: IProps) => (
  <Form.Item className={`input-headers ${className}`} label={`${label} list`} shouldUpdate>
    <Form.List name={name.length === 1 ? name[0] : name} initialValue={initialValue}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.HeaderContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder={`${unit} ${index + 1}`} />
              </Form.Item>

              <Form.Item name={[field.name, 'value']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder={`Value ${index + 1}`} />
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
            Add {label}
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default RequestDetailsHeadersInput;
