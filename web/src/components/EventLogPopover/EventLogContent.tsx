import {useEffect, useRef} from 'react';
import TestRunEvent from 'models/TestRunEvent.model';
import EventLogService from 'services/EventLog.service';
import * as S from './EventLogPopover.styled';

interface IProps {
  runEvents: TestRunEvent[];
}

const EventLogContent = ({runEvents}: IProps) => {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({behavior: 'smooth'});
  }, [runEvents]);

  return (
    <S.Container>
      {runEvents.map(event => (
        <S.EventEntry key={`${event.type}-${event.createdAt}`} $logLevel={event.logLevel}>
          <b>{EventLogService.typeToString(event)}</b> {EventLogService.detailsToString(event)}
        </S.EventEntry>
      ))}
      <div ref={bottomRef} />
    </S.Container>
  );
};

export default EventLogContent;
