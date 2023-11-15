import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor} from 'components/Inputs';
import * as S from './MultiURL.styled';

interface IProps {
  name?: string[];
}

function isFirstItem(index: number) {
  return index === 0;
}

const MultiURL = ({name = ['brokerUrls']}: IProps) => (
  <Form.Item
    rules={[{required: true, min: 1, message: 'Please enter a valid URL'}]}
    shouldUpdate
    style={{marginBottom: 0}}
  >
    <Form.List name={name.length === 1 ? name[0] : name} initialValue={['']}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.URLContainer key={field.name}>
              <Form.Item name={[field.name]} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder={`Enter broker URL (${index + 1})`} />
              </Form.Item>

              {!isFirstItem(index) && (
                <Form.Item noStyle>
                  <Button icon={<S.DeleteIcon />} onClick={() => remove(field.name)} type="link" />
                </Form.Item>
              )}
            </S.URLContainer>
          ))}

          <S.AddURLContainer>
            <Button
              icon={<PlusOutlined />}
              onClick={() => add()}
              style={{fontWeight: 600, height: 'auto', padding: 0}}
              type="link"
            >
              New URL
            </Button>
          </S.AddURLContainer>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default MultiURL;
