import {Typography} from 'antd';
import {capitalize} from 'lodash';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {IPropsEvent} from './RunEvent';
import * as S from './RunEvents.styled';

const RunEventPolling = ({event}: IPropsEvent) => {
  const {
    pollingProfile: {
      strategy,
      periodic: {retryDelay = '', timeout = ''},
    },
  } = useSettingsValues();

  return (
    <S.EventContainer>
      <S.Dot $logLevel={event.logLevel} />
      <Typography.Title level={3}>{event.title}</Typography.Title>
      <Typography.Text>{event.description}</Typography.Text>
      {!!event.polling && (
        <S.Content>
          <S.Column>
            <S.InfoIcon />
            <div>
              <Typography.Title level={3}>Polling profile configuration:</Typography.Title>
              <S.Column>
                <Typography.Text>{capitalize(strategy)} strategy</Typography.Text>
              </S.Column>
              <S.Column>
                <Typography.Text>{retryDelay} of retry delay</Typography.Text>
              </S.Column>
              <S.Column>
                <Typography.Text>{timeout} of timeout</Typography.Text>
              </S.Column>
            </div>
          </S.Column>
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
        </S.Content>
      )}
    </S.EventContainer>
  );
};

export default RunEventPolling;
