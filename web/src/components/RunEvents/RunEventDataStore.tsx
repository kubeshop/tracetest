import {Typography} from 'antd';
import {Link} from 'react-router-dom';
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
        {event.dataStoreConnection?.allPassed ? (
          <Typography.Title level={3}>Data store configuration is valid.</Typography.Title>
        ) : (
          <S.InvalidDataStoreContainer>
            <Typography.Title level={3} style={{margin: 0}}>
              Data store configuration is not valid.
            </Typography.Title>
            <Typography.Text>
              You can go to the <Link to="/settings">settings page</Link> to fix it.
            </Typography.Text>
          </S.InvalidDataStoreContainer>
        )}

        <TestConnectionNotification result={event.dataStoreConnection} />
      </S.Content>
    )}
  </S.EventContainer>
);

export default RunEventDataStore;
