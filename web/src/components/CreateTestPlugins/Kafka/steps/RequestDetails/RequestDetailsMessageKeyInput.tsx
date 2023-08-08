import {Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

interface IProps {

}

const RequestDetailsMessageKey = (form: IProps) => {
  return (
    <div>
      <S.Label>Message Key</S.Label>
      <Form.Item data-cy="message-key" name="message-key" rules={[{required: true, message: 'Please enter a message key'}]}>
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter a message key / identifier (Optional)" />
      </Form.Item>
    </div>
  );
};

export default RequestDetailsMessageKey;
