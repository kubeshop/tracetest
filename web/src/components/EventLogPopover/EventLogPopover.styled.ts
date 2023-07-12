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
      padding: 5px 16px;
    }
  }
`;

export const Container = styled.div`
  margin: 0;
  max-height: calc(100vh - 200px);
  width: 550px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 2px;
  border-radius: 2px;
  position: relative;
`;

export const CopyIcon = styled(CopyOutlined)`
  color: ${({theme}) => theme.color.primary};
  cursor: pointer;
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

export const TitleContainer = styled.div`
  padding: 5px 0;
  display: flex;
  justify-content: space-between;
`;

export const Title = styled(Typography.Title)`
  && {
    margin: 0;
    font-size: ${({theme}) => theme.size.md};
  }
`;
