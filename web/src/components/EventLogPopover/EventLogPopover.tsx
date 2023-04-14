import {Popover} from 'antd';
import TestRunEvent from 'models/TestRunEvent.model';
import EventLogContent from './EventLogContent';
import * as S from './EventLogPopover.styled';

interface IProps {
  runEvents: TestRunEvent[];
}

const EventLogPopover = ({runEvents}: IProps) => {
  return (
    <>
      <S.GlobalStyle />
      <Popover
        id="eventlog-popover"
        content={<EventLogContent runEvents={runEvents} />}
        trigger="click"
        placement="bottomLeft"
      >
        <S.TerminalIcon />
      </Popover>
    </>
  );
};

export default EventLogPopover;
