import {Typography} from 'antd';
import {TConnectionTestStep} from 'types/Config.types';
import * as S from './TestConnectionNotification.styled';

interface IProps {
  step: TConnectionTestStep;
}

const TestConnectionStep = ({step: {message, error: err, passed}}: IProps) => {
  return message || err ? (
    <S.StepContainer>
      {passed ? <S.SuccessCheckIcon /> : <S.FailedCheckIcon />}
      <Typography.Text>{message || err}</Typography.Text>
    </S.StepContainer>
  ) : null;
};

export default TestConnectionStep;
