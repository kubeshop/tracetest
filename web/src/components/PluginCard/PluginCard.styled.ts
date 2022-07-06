import {Typography} from 'antd';
import styled from 'styled-components';

export const PluginCard = styled.div<{$isSelected: boolean; $isActive: boolean}>`
  display: flex;
  border-radius: 4px;
  border: 1px solid ${({$isSelected, theme}) => ($isSelected ? theme.color.primary : theme.color.borderLight)};
  padding: 20px;
  background: ${({$isSelected, theme}) => ($isSelected ? theme.color.background : theme.color.white)};
  gap: 20px;
  opacity: ${({$isActive}) => ($isActive ? 1 : 0.5)};
  min-width: 490px;
  cursor: ${({$isActive}) => ($isActive ? 'pointer' : 'default')};
`;

export const Content = styled.div`
  display: flex;
  flex-direction: column;
`;

export const Title = styled(Typography.Text).attrs({
  strong: true,
})``;
export const Description = styled(Typography.Text)``;

export const Circle = styled.div`
  margin-top: 4px;
  height: 16px;
  width: 16px;
  margin-left: 5px;
  border: ${({theme}) => `1px solid ${theme.color.primary}`};
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
`;

export const Check = styled.div`
  height: 8px;
  width: 8px;
  background: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  display: inline-block;
`;
