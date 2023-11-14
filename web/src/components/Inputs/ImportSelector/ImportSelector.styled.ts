import styled from 'styled-components';
import {Typography} from 'antd';

export const CardContainer = styled.div`
  align-items: center;
  background: ${({theme}) => theme.color.white};
  border: 1px solid ${({theme}) => theme.color.borderLight};
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  gap: 12px;
  padding: 16px 12px;
  width: 48%;

  &:hover {
    background: ${({theme}) => theme.color.background};
    border: 1px solid ${({theme}) => theme.color.primary};
  }
`;

export const CardContent = styled.div`
  display: flex;
  flex-direction: column;
`;

export const CardDescription = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.xs};
  opacity: 1;
`;

export const CardList = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
`;

export const CardTitle = styled(Typography.Text).attrs({
  strong: true,
})`
  display: inline-block;
  font-size: ${({theme}) => theme.size.sm};
  opacity: 1;

  a {
    opacity: 1;
  }
`;

export const Check = styled.div`
  background: ${({theme}) => theme.color.primary};
  border-radius: 50%;
  display: inline-block;
  height: 8px;
  width: 8px;
`;

export const Circle = styled.div`
  align-items: center;
  border-radius: 50%;
  border: ${({theme}) => `1px solid ${theme.color.primary}`};
  display: flex;
  justify-content: center;
  max-height: 16px;
  max-width: 16px;
  min-height: 16px;
  min-width: 16px;
  opacity: 1;
`;
