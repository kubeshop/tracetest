import {Typography} from 'antd';
import TestConnectionNotification from 'components/TestConnectionNotification/TestConnectionNotification';
import {IPropsEvent} from './RunEvent';
import * as S from './RunEvents.styled';

const RunEventDataStore = ({event}: IPropsEvent) => (
  <S.EventContainer>
    <S.Dot $logLevel={event.logLevel} />
    <Typography.Title level={3}>{event.title}</Typography.Title>
    <Typography.Text>{event.description}</Typography.Text>
    {!!event.dataStoreConnection && (
      <S.Content>
        <Typography.Title level={3}>
          {event.dataStoreConnection?.allPassed
            ? 'Data store configuration is valid.'
            : 'Data store configuration is not valid.'}
        </Typography.Title>
        <TestConnectionNotification result={event.dataStoreConnection} />
      </S.Content>
    )}
  </S.EventContainer>
);

export default RunEventDataStore;
