import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import {IKeyValue} from 'constants/Test.constants';
import * as S from './KeyValueList.styled';
import SingleLine from '../../Inputs/SingleLine';

interface IProps {
  name?: string;
  className?: string;
  label?: string;
  addButtonLabel?: string;
  keyPlaceholder?: string;
  valuePlaceholder?: string;
  initialValue?: IKeyValue[];
}
const KeyValueList = ({
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
              <S.Item>
                <Form.Item name={[field.name, 'key']} noStyle>
                  <SingleLine placeholder={`${keyPlaceholder} ${index + 1}`} />
                </Form.Item>
              </S.Item>

              <S.Item>
                <Form.Item name={[field.name, 'value']} noStyle>
                  <SingleLine placeholder={`${valuePlaceholder} ${index + 1}`} />
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

export default KeyValueList;
