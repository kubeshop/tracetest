import {Typography} from 'antd';
import TestRunEvent from 'models/TestRunEvent.model';
import * as S from './RunEvents.styled';

export interface IPropsEvent {
  event: TestRunEvent;
}

const RunEvent = ({event}: IPropsEvent) => (
  <S.EventContainer>
    <S.Dot $logLevel={event.logLevel} />
    <Typography.Title level={3}>{event.title}</Typography.Title>
    <Typography.Text>{event.description}</Typography.Text>
  </S.EventContainer>
);

export default RunEvent;
