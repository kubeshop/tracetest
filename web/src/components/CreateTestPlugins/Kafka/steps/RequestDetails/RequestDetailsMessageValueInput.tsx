import {Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

interface IProps {

}

const RequestDetailsMessageValue = (form: IProps) => {
  return (
    <div>
      <S.Label>Message Value</S.Label>
      <Form.Item data-cy="message-value" name="message-value" rules={[{required: true, message: 'Please enter a message value'}]}>
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter a message value / content" />
      </Form.Item>
    </div>
  );
};

export default RequestDetailsMessageValue;
