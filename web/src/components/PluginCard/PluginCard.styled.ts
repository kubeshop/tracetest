import {Typography} from 'antd';
import styled from 'styled-components';

export const PluginCard = styled.div<{$isSelected: boolean; $isActive: boolean}>`
  display: flex;
  border-radius: 4px;
  border: 1px solid ${({$isSelected, theme}) => ($isSelected ? theme.color.primary : theme.color.borderLight)};
  padding: 4px;
  padding-left: 16px;
  background: ${({$isSelected, theme}) => ($isSelected ? theme.color.background : theme.color.white)};
  gap: 12px;
  width: 48%;
  cursor: ${({$isActive}) => ($isActive ? 'pointer' : 'default')};
  align-items: center;
`;

export const Content = styled.div`
  display: flex;
  flex-direction: column;
`;

export const Title = styled(Typography.Text).attrs({
  strong: true,
})<{$isActive: boolean}>`
  display: inline-block;
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};
  font-size: ${({theme}) => theme.size.sm};

  a {
    opacity: 1;
  }
`;
export const Description = styled(Typography.Text)<{$isActive: boolean}>`
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};
  font-size: ${({theme}) => theme.size.xs};
`;

export const Circle = styled.div<{$isActive: boolean}>`
  min-height: 16px;
  min-width: 16px;
  max-height: 16px;
  max-width: 16px;
  border: ${({theme}) => `1px solid ${theme.color.primary}`};
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};
`;

export const Check = styled.div`
  height: 8px;
  width: 8px;
  background: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  display: inline-block;
`;
