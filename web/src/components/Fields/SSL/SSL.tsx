import {Form, Switch} from 'antd';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import * as S from './SSL.styled';

interface IProps {
  formID?: string;
}

const SSL = ({formID}: IProps) => (
  <S.SSLVerificationContainer>
    <label htmlFor={`${formID}_sslVerification`}>Enable SSL certificate verification</label>
    <Form.Item name="sslVerification" valuePropName="checked" style={{marginBottom: 0}}>
      <Switch />
    </Form.Item>
    <TooltipQuestion title="Verify SSL certificates when sending the request. Verification failures will result in the request being aborted." />
  </S.SSLVerificationContainer>
);

export default SSL;
