import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  background: ${({theme}) => theme.color.background};
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  display: grid;
  gap: 16px;
  grid-template-columns: 1fr 80%;
  margin-bottom: 8px;
  min-height: 58px;
  padding: 7px 16px;
`;

export const Column = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

export const Text = styled(Typography.Text)<{$fontWeight?: string; $opacity?: number}>`
  && {
    color: ${({$opacity}) => `rgba(3, 24, 73, ${$opacity || 1})`};
    font-size: ${({theme}) => theme.size.sm};
    font-weight: ${({$fontWeight}) => $fontWeight || 'normal'};
    margin: 0;
  }
`;
