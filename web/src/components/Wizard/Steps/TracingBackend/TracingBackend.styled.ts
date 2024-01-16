import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div``;

export const BackendSelector = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 16px;
`;

export const Card = styled.div`
  display: flex;
  gap: 6px;
  align-items: center;
  padding: 6px;
  border: 1px solid ${({theme}) => theme.color.border};
  border-radius: 2px;
  cursor: pointer;
  height: 32px;
  width: 142px;
  font-size: ${({theme}) => theme.size.sm};
  font-weight: 700;

  &:hover {
    border-color: ${({theme}) => theme.color.primary};
    background-color: ${({theme}) => theme.color.background};
  }
`;

export const Name = styled(Typography.Text)`
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
`;

export const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
`;