import {Popover, Tooltip} from 'antd';
import {useState} from 'react';
import {CloseOutlined, ForwardOutlined} from '@ant-design/icons';
import * as S from './SkipPollingPopover.styled';
import Content from './Content';

interface IProps {
  isLoading: boolean;
  skipPolling(shouldSave: boolean): void;
}

const SkipPollingPopover = ({isLoading, skipPolling}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <S.StopContainer>
      <S.GlobalStyle />
      <Popover
        id="skip-trace-popover"
        title={
          <S.ContentContainer>
            <S.Title level={3}>Waiting too long for the trace?</S.Title>
            <S.CloseIcon onClick={() => setIsOpen(false)}>
              <CloseOutlined />
            </S.CloseIcon>
          </S.ContentContainer>
        }
        content={<Content isLoading={isLoading} skipPolling={skipPolling} />}
        visible={isOpen}
        placement="bottomRight"
      >
        <Tooltip title="Skip Trace collection" placement="left">
          <S.ForwardButton size="small" loading={isLoading} onClick={() => setIsOpen(true)} type="link">
            <ForwardOutlined />
          </S.ForwardButton>
        </Tooltip>
      </Popover>
    </S.StopContainer>
  );
};

export default SkipPollingPopover;
