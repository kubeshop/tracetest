import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import {IKeyValue} from 'constants/Test.constants';
import * as S from './KeyValueListInput.styled';

interface IProps {
  name?: string;
  className?: string;
  label?: string;
  addButtonLabel?: string;
  keyPlaceholder?: string;
  valuePlaceholder?: string;
  initialValue?: IKeyValue[];
}
const KeyValueInputList = ({
  name = 'headers',
  className = '',
  label = '',
  addButtonLabel = '',
  keyPlaceholder = '',
  valuePlaceholder = '',
  initialValue = [],
}: IProps) => (
  <Form.Item className={`input-headers ${className}`} label={label} shouldUpdate>
    <Form.List name={name} initialValue={initialValue}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.KeyValueContainer key={field.name}>
              <Form.Item name={[field.name, 'key']} style={{marginBottom: 0, marginRight: '2px'}}>
                <Editor type={SupportedEditors.Interpolation} placeholder={`${keyPlaceholder} ${index + 1}`} />
              </Form.Item>

              <Form.Item name={[field.name, 'value']} style={{marginBottom: 0}}>
                <Editor type={SupportedEditors.Interpolation} placeholder={`${valuePlaceholder} ${index + 1}`} />
              </Form.Item>

              <Form.Item noStyle>
                <Button
                  icon={<S.DeleteIcon />}
                  onClick={() => remove(field.name)}
                  style={{marginLeft: 12}}
                  type="link"
                />
              </Form.Item>
            </S.KeyValueContainer>
          ))}

          <Button
            data-cy="add-header"
            icon={<PlusOutlined />}
            onClick={() => add()}
            style={{fontWeight: 600, height: 'auto', padding: 0}}
            type="link"
          >
            {addButtonLabel}
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default KeyValueInputList;
