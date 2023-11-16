import {Button, Popover} from 'antd';
import {differenceInSeconds} from 'date-fns';
import {useEffect, useState} from 'react';
import {ForwardOutlined} from '@ant-design/icons';
import * as S from './SkipPollingPopover.styled';
import Content from './Content';

interface IProps {
  isLoading: boolean;
  skipPolling(shouldSave: boolean): void;
  startTime: string;
}

const TIMEOUT_TO_SHOW = 10; // seconds

const SkipPollingPopover = ({isLoading, skipPolling, startTime}: IProps) => {
  const [isOpen, setIsOpen] = useState(true);
  const diff = differenceInSeconds(new Date(), new Date(startTime));

  useEffect(() => {
    if (diff > TIMEOUT_TO_SHOW) setIsOpen(true);
  }, [diff, isOpen]);

  return (
    <S.StopContainer>
      <S.GlobalStyle />
      <Popover
        id="skip-trace-popover"
        title={<S.Title level={3}>Taking too long to get the trace?</S.Title>}
        content={<Content isLoading={isLoading} skipPolling={skipPolling} />}
        visible={isOpen}
        placement="bottomRight"
      >
        <Button loading={isLoading} onClick={() => skipPolling(false)} type="link">
          <ForwardOutlined />
        </Button>
      </Popover>
    </S.StopContainer>
  );
};

export default SkipPollingPopover;
