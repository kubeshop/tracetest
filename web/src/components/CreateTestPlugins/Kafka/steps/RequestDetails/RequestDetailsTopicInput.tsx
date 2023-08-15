import {Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

const RequestDetailsTopic = () => {
  return (
    <div>
      <S.Label>Topic</S.Label>
      <Form.Item data-cy="topic" name="topic" rules={[{required: true, message: 'Please enter a topic'}]}>
        <Editor type={SupportedEditors.Interpolation} placeholder="Enter a topic / message category" />
      </Form.Item>
    </div>
  );
};

export default RequestDetailsTopic;
