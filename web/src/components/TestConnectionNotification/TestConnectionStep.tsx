import {Typography} from 'antd';
import ConnectionTestStep from 'models/ConnectionResultStep.model';
import * as S from './TestConnectionNotification.styled';

interface IProps {
  step: ConnectionTestStep;
  title: string;
}

const iconMap = {
  passed: <S.SuccessCheckIcon />,
  failed: <S.FailedCheckIcon />,
  warning: <S.WarningCheckIcon />,
};

const TestConnectionStep = ({step: {message, error: err, status}, title}: IProps) => {
  const icon = iconMap[status];

  return message || err ? (
    <S.StepContainer>
      {icon}
      <div>
        <S.Title level={3}>{title}</S.Title>
        <Typography.Text>{message}</Typography.Text>
        {!!err && <Typography.Text type="secondary"> - Error: {err}</Typography.Text>}
      </div>
    </S.StepContainer>
  ) : null;
};

export default TestConnectionStep;
