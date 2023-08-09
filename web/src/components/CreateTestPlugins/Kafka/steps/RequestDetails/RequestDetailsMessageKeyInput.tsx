import {Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

const RequestDetailsMessageKey = () => {
  return (
    <div>
      <S.Label>Message Key</S.Label>
      <Form.Item data-cy="message-key" name="messageKey">
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter a message key / identifier (Optional)" />
      </Form.Item>
    </div>
  );
};

export default RequestDetailsMessageKey;
