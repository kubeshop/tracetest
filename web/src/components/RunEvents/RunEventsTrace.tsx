import {Typography} from 'antd';

import {TRACE_DOCUMENTATION_URL} from 'constants/Common.constants';
import {TestState} from 'constants/TestRun.constants';
import {TraceEventType} from 'constants/TestRunEvents.constants';
import TestRunService from 'services/TestRun.service';
import RunEvent, {IPropsEvent} from './RunEvent';
import RunEventDataStore from './RunEventDataStore';
import RunEventPolling from './RunEventPolling';
import {IPropsComponent} from './RunEvents';
import * as S from './RunEvents.styled';

const ComponentMap: Record<TraceEventType, (props: IPropsEvent) => React.ReactElement> = {
  [TraceEventType.DATA_STORE_CONNECTION_INFO]: RunEventDataStore,
  [TraceEventType.POLLING_ITERATION_INFO]: RunEventPolling,
};

const RunEventsTrace = ({events, state}: IPropsComponent) => {
  const filteredEvents = TestRunService.getTestRunEventsWithLastPolling(events);

  const loadingHeader = (
    <>
      <S.LoadingIcon />
      <Typography.Title level={3} type="secondary">
        We are working to gather the trace associated with this test run
      </Typography.Title>
      <S.Paragraph type="secondary">
        Want to know more about traces? Head to the official{' '}
        <S.Link href={TRACE_DOCUMENTATION_URL} target="_blank">
          Open Telemetry Documentation
        </S.Link>
      </S.Paragraph>
    </>
  );

  const failedHeader = (
    <>
      <S.ErrorIcon />
      <Typography.Title level={2} type="secondary">
        Trace Fetch Failed
      </Typography.Title>
      <S.Paragraph type="secondary">
        We prepared the breakdown of diagnostic steps and tips to help identify the source of the issue:
      </S.Paragraph>
    </>
  );

  return (
    <S.Container $hasScroll>
      {state === TestState.FAILED ? failedHeader : loadingHeader}

      <S.ListContainer>
        {filteredEvents.map(event => {
          const Component = ComponentMap[event.type as TraceEventType] ?? RunEvent;
          return <Component event={event} key={event.type} />;
        })}
      </S.ListContainer>
    </S.Container>
  );
};

export default RunEventsTrace;
