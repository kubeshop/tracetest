import {PlusOutlined} from '@ant-design/icons';
import {Button, Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

interface IProps {
  name?: string[];
}

function isFirstItem(index: number) {
  return index === 0;
}

const RequestDetailsUrlInput = ({
  name = ['brokerUrls'],
}: IProps) => (
  <Form.Item className="input-url" label="Broker URLs" shouldUpdate>
    <Form.List name={name.length === 1 ? name[0] : name}>
      {(fields, {add, remove}) => (
        <>
          {fields.map((field, index) => (
            <S.BrokerURLInputContainer $firstItem={isFirstItem(index)} key={field.name}>
              <Form.Item name={[field.name, 'url']} noStyle>
                <Editor type={SupportedEditors.Interpolation} placeholder={`Enter a broker URL (${index + 1})`} />
              </Form.Item>

              {!isFirstItem(index) && <Form.Item noStyle>
                <Button
                  icon={<S.DeleteIcon />}
                  onClick={() => remove(field.name)}
                  style={{marginLeft: 12}}
                  type="link"
                />
              </Form.Item>}
            </S.BrokerURLInputContainer>
          ))}

          <Button
            data-cy="add-broker-url"
            icon={<PlusOutlined />}
            onClick={() => add()}
            style={{fontWeight: 600, height: 'auto', padding: 0}}
            type="link"
          >
            Add URL
          </Button>
        </>
      )}
    </Form.List>
  </Form.Item>
);

export default RequestDetailsUrlInput;
