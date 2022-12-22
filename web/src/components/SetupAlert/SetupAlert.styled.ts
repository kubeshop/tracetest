import styled from 'styled-components';
import {Alert, Button, Typography} from 'antd';

export const Container = styled(Alert)`
  position: absolute;
  z-index: 101;
  width: 100%;
`;

export const TextBold = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    font-weight: 600;
    white-space: nowrap;
  }
`;

export const Text = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.md};
    font-weight: 400;
  }
`;

export const Message = styled.div`
  display: flex;
  gap: 12px;
  align-items: center;
  margin-left: 16px;
`;

export const WarningButton = styled(Button)`
  background: ${({theme}) => theme.color.warningYellow};
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border: none;
  font-weight: 600;

  &:hover,
  &:focus {
    background: ${({theme}) => theme.color.warningYellow};
  }
`;
