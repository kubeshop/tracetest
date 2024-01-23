import OTLPTestConnectionResponse from 'models/OTLPTestConnectionResponse.model';
import Date from 'utils/Date';
import {Button} from 'antd';
import * as S from './TestConnectionStatus.styled';
import useTestConnectionStatus, {TConnectionStatus} from './hooks/useTestConnnectionStatus';

const Icon = ({status}: {status: TConnectionStatus}) => {
  switch (status) {
    case 'loading':
      return <S.LoadingIcon />;
    case 'success':
      return <S.IconSuccess />;
    case 'error':
      return <S.IconFail />;
    case 'idle':
    default:
      return <S.IconInfo />;
  }
};

const getText = (status: TConnectionStatus, {spanCount, lastSpanTimestamp} = OTLPTestConnectionResponse()) => {
  switch (status) {
    case 'loading':
      return 'Waiting for incoming traces';
    case 'success':
      return (
        <>
          We received your traces
          <br />
          <S.SpanCountTest>
            <i>
              {spanCount} spans, {Date.getTimeAgo(lastSpanTimestamp)}
            </i>
          </S.SpanCountTest>
        </>
      );
    case 'error':
      return 'We could not receive your traces';
    case 'idle':
    default:
      return 'Waiting for Test Connection';
  }
};

interface IProps {
  onTestConnection(): void;
}

const TestConnectionStatus = ({onTestConnection}: IProps) => {
  const {status, isOtlpBased, isLoading, otlpResponse} = useTestConnectionStatus();

  if (!isOtlpBased)
    return (
      <Button loading={isLoading} type="primary" ghost onClick={onTestConnection}>
        Test Connection
      </Button>
    );

  return (
    <S.StatusContainer>
      <Icon status={status} />
      <S.StatusText>{getText(status, otlpResponse)}</S.StatusText>
    </S.StatusContainer>
  );
};

export default TestConnectionStatus;
