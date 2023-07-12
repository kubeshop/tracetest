import {Popover, Tooltip} from 'antd';
import TestRunEvent from 'models/TestRunEvent.model';
import useCopy from 'hooks/useCopy';
import EventLogService from 'services/EventLog.service';
import EventLogContent from './EventLogContent';
import * as S from './EventLogPopover.styled';

interface IProps {
  runEvents: TestRunEvent[];
}

const EventLogPopover = ({runEvents}: IProps) => {
  const copy = useCopy();

  return (
    <>
      <S.GlobalStyle />
      <Popover
        id="eventlog-popover"
        content={<EventLogContent runEvents={runEvents} />}
        trigger="click"
        placement="bottomLeft"
        title={
          <S.TitleContainer>
            <S.Title>Event Log</S.Title>
            <Tooltip title="Copy Text">
              <S.CopyIcon onClick={() => copy(EventLogService.listToString(runEvents))} />
            </Tooltip>
          </S.TitleContainer>
        }
      >
        <S.TerminalIcon />
      </Popover>
    </>
  );
};

export default EventLogPopover;
