import {Form} from 'antd';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from 'components/Editor';
import * as S from './RequestDetails.styled';

interface IProps {

}

const RequestDetailsUrlInput = (props: IProps) => {
  return (
    <div>
      <S.Label>Broker URLs</S.Label>
      <S.BrokerURLInputContainer>
        <Form.Item data-cy="broker-url" name="broker-url" rules={[{required: true, message: 'Please enter a broker url'}]}>
          <Editor type={SupportedEditors.Interpolation} placeholder="Enter broker url" />
        </Form.Item>
      </S.BrokerURLInputContainer>
    </div>
  );

  // TODO daniel: add button
};

export default RequestDetailsUrlInput;
