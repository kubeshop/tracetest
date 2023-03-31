import {CloseCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled, {DefaultTheme} from 'styled-components';

import {LogLevel} from 'models/TestRunEvent.model';

function getLogLevelColor(logLevel: LogLevel, theme: DefaultTheme): string {
  if (logLevel === LogLevel.Error) return theme.color.error;
  if (logLevel === LogLevel.Success) return theme.color.success;
  if (logLevel === LogLevel.Warning) return theme.color.warningYellow;
  return theme.color.textLight;
}

export const Container = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 70px 120px;
`;

export const Dot = styled.div<{$logLevel?: LogLevel}>`
  background: #ffffff;
  border: 2px solid ${({theme, $logLevel = LogLevel.Info}) => getLogLevelColor($logLevel, theme)};
  border-radius: 50%;
  left: -5px;
  height: 10px;
  position: absolute;
  top: 0;
  width: 10px;
`;

export const ErrorIcon = styled(CloseCircleFilled)`
  color: ${({theme}) => theme.color.error};
  font-size: 32px;
  margin-bottom: 26px;
`;

export const EventContainer = styled.div`
  border-left: 1px solid ${({theme}) => theme.color.textLight};
  padding: 0 20px 20px;
  position: relative;
`;

export const Link = styled(Typography.Link)`
  font-weight: bold;
`;

export const ListContainer = styled.div`
  padding: 24px 0;
`;

export const Paragraph = styled(Typography.Paragraph)`
  text-align: center;
`;
