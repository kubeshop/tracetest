import {SettingOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import styled from 'styled-components';

export const Icon = styled(SettingOutlined)`
  && {
    cursor: pointer;
    font-size: ${({theme}) => theme.size.lg};
`;

export const WarningButton = styled(Button)`
  background: ${({theme}) => theme.color.warningYellow};
  width: 57px;
  height: 24px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 0 !important;
  border: none;
  font-weight: 600;

  &:hover,
  &:focus {
    background: ${({theme}) => theme.color.warningYellow};
  }
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin: 0;
  }
`;

export const Text = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.md};
    margin: 0;
  }
`;

export const MessageContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-width: 340px;
`;

export const ButtonContainer = styled.div`
  justify-content: flex-end;
  display: flex;
  margin-top: 6px;
`;

export const Trigger = styled.div`
  display: flex;
  gap: 6px;
  align-items: center;
  cursor: pointer;
  margin-right: 8px;
  border-radius: 12px;
  padding: 4px 7px;
  background: ${({theme}) => theme.color.backgroundDark};
  height: 24px;
`;
