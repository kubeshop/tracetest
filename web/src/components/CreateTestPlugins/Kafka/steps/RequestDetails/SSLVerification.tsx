import {Form, Switch} from 'antd';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import * as S from './RequestDetails.styled';

const SSLVerification = () => {
  return (
    <S.SSLVerificationContainer>
      <label htmlFor="sslVerification">Enable SSL certificate verification</label>
      <Form.Item name="sslVerification" valuePropName="checked">
        <Switch />
      </Form.Item>
      <TooltipQuestion title="Verify SSL certificates when connection to broker. Verification failures will result in the broker connection being aborted." />
    </S.SSLVerificationContainer>
  );
};

export default SSLVerification;
