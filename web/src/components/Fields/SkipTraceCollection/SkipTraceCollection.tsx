import {Form, Switch} from 'antd';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import * as S from '../SSL/SSL.styled';

const SkipTraceCollection = () => {
  return (
    <S.SSLVerificationContainer>
      <label htmlFor="skipTraceCollection">Skip Trace Collection</label>
      <Form.Item name="skipTraceCollection" valuePropName="checked" style={{marginBottom: 0}}>
        <Switch id="skipTraceCollection" />
      </Form.Item>
      <TooltipQuestion title="Skips Trace Collection for all runs. You can still create and run tests." />
    </S.SSLVerificationContainer>
  );
};

export default SkipTraceCollection;
