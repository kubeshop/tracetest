import {Typography} from 'antd';
import ConnectionTestStep from 'models/ConnectionResultStep.model';
import * as S from './TestConnectionNotification.styled';

interface IProps {
  step: ConnectionTestStep;
  title: string;
}

const TestConnectionStep = ({step: {message, error: err, passed}, title}: IProps) => {
  return message || err ? (
    <S.StepContainer>
      {passed ? <S.SuccessCheckIcon /> : <S.FailedCheckIcon />}
      <div>
        <S.Title level={3}>{title}</S.Title>
        <Typography.Text>{message}</Typography.Text>
        {!!err && <Typography.Text type="secondary"> - Error: {err}</Typography.Text>}
      </div>
    </S.StepContainer>
  ) : null;
};

export default TestConnectionStep;
