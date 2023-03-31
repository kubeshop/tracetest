import {Typography} from 'antd';
import {IPropsEvent} from './RunEvent';
import * as S from './RunEvents.styled';

const RunEventPolling = ({event}: IPropsEvent) => (
  <S.EventContainer>
    <S.Dot $logLevel={event.logLevel} />
    <Typography.Title level={3}>{event.title}</Typography.Title>
    <Typography.Text>{event.description}</Typography.Text>
    {!!event.polling && (
      <S.Content>
        <S.Column>
          <S.InfoIcon />
          <div>
            <Typography.Title level={3}>Number of spans collected:</Typography.Title>
            <Typography.Text>{event.polling.periodic.numberSpans}</Typography.Text>
          </div>
        </S.Column>
        <S.Column>
          <S.InfoIcon />
          <div>
            <Typography.Title level={3}>Iteration number:</Typography.Title>
            <Typography.Text>{event.polling.periodic.numberIterations}</Typography.Text>
          </div>
        </S.Column>
        <S.Column>
          <S.InfoIcon />
          <div>
            <Typography.Title level={3}>Reason why the next iteration will be executed:</Typography.Title>
            <Typography.Text>{event.polling.reasonNextIteration}</Typography.Text>
          </div>
        </S.Column>
      </S.Content>
    )}
  </S.EventContainer>
);

export default RunEventPolling;
