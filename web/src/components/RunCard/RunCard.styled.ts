import {Typography} from 'antd';
import styled from 'styled-components';

export const ResultCard = styled.div`
  display: grid;
  align-items: center;
  grid-template-columns: 300px 300px 100px 220px 40px 40px 40px 1fr;
  gap: 16px;
  padding: 16px 12px;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  background: ${({theme}) => theme.color.background};
  cursor: pointer;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
  overflow-x: ellipsis;
`;
