import {Button, Space, Typography} from 'antd';
import styled from 'styled-components';

import emptyStateIcon from 'assets/SpanAssertionsEmptyState.svg';

export const Container = styled.div`
  background: ${({theme}) => theme.color.white};
  border-left: 4px solid #2f1e61;
  border-radius: 2px;
  box-shadow: -1px 1px 5px #e4e9f5;
`;

export const Box = styled.div`
  align-items: center;
  background: #2f1e61;
  border-radius: 3px;
  display: flex;
  height: 38px;
  justify-content: center;
  margin-left: 26px;
  margin-right: 18px;
  width: 38px;
`;

export const BoxTitle = styled(Typography.Title)`
  && {
    color: ${({theme}) => theme.color.white};
    margin: 0;
  }
`;

export const EmptyStateIcon = styled.img.attrs({
  src: emptyStateIcon,
})``;

export const EmptyStateContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 14px;
  justify-content: center;
  padding: 16px 0;
`;

export const HeaderDetail = styled(Typography.Text)`
  display: flex;
  align-items: center;
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
  margin-right: 8px;
`;

export const Link = styled(Button).attrs({
  type: 'link',
})`
  color: ${({theme}) => theme.color.primary};
  font-weight: 600;
  padding: 0;
`;

export const LoadingContainer = styled(Space)`
  width: 100%;
`;

export const HeaderDot = styled.span<{$passed: boolean}>`
  background-color: ${({$passed, theme}) => ($passed ? theme.color.success : theme.color.error)};
  border-radius: 50%;
  display: inline-block;
  height: 10px;
  line-height: 0;
  margin-right: 4px;
  vertical-align: -0.1em;
  width: 10px;
`;

export const FooterContainer = styled.div`
  margin-top: 8px;
  text-align: right;
  width: 100%;
`;

export const Row = styled.div<{$gap?: number}>`
  align-items: center;
  column-gap: ${({$gap}) => $gap && `${$gap}px`};
  display: flex;
`;

export const RunButton = styled(Button)`
  margin-right: 12px;
`;

export const RunsContainer = styled.div`
  padding: 0 24px 15px 64px;
`;

export const TestContainer = styled.div`
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  padding: 15px 24px;
`;

export const Text = styled(Typography.Text).attrs({as: 'p'})`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
  margin: 0;
`;

export const Title = styled(Typography.Title)`
  && {
    margin: 0;
  }
`;
