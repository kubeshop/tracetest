import {LeftOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const BackIcon = styled(LeftOutlined)`
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const Container = styled.div<{$isWhite?: boolean}>`
  align-items: center;
  background: ${({$isWhite, theme}) => ($isWhite ? theme.color.white : theme.color.background)};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  width: 100%;
  margin-bottom: 16px;
`;

export const Section = styled.div`
  align-items: center;
  display: flex;
  gap: 14px;
`;

export const Text = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  && {
    font-size: ${({theme}) => theme.size.sm};
    margin: 0;
  }
`;

export const Title = styled(Typography.Title).attrs({level: 2})`
  && {
    margin: 0;
  }
`;
