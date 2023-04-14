import {CopyOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled, {DefaultTheme, createGlobalStyle} from 'styled-components';
import {LogLevel} from 'constants/TestRunEvents.constants';
import terminalIcon from 'assets/terminal.svg';

function getLogLevelColor(logLevel: LogLevel, theme: DefaultTheme): string {
  if (logLevel === LogLevel.Error) return theme.color.error;
  if (logLevel === LogLevel.Success) return theme.color.success;
  if (logLevel === LogLevel.Warning) return theme.color.warningYellow;
  return theme.color.text;
}

export const GlobalStyle = createGlobalStyle`
  #eventlog-popover {
    .ant-popover-inner-content {
      padding: 5px;
    }
  }
`;

export const Container = styled.div`
  margin: 0;
  max-height: calc(100vh - 200px);
  width: 550px;
  overflow-y: auto;
  padding: 5px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-radius: 2px;
  position: relative;
`;

export const CopyIconContainer = styled.div`
  position: absolute;
  right: 8px;
  top: 9px;
  padding: 0 2px;
  border-radius: 2px;
  cursor: pointer;
  z-index: 101;
`;

export const CopyIcon = styled(CopyOutlined)`
  color: ${({theme}) => theme.color.text};
`;

export const EventEntry = styled(Typography.Text)<{$logLevel?: LogLevel}>`
  font-size: ${({theme}) => theme.size.sm};
  color: ${({theme, $logLevel = LogLevel.Info}) => getLogLevelColor($logLevel, theme)};
  margin: 0;
`;

export const TerminalIcon = styled.img.attrs({src: terminalIcon})`
  width: 20px;
  height: 20px;
  cursor: pointer;
`;
