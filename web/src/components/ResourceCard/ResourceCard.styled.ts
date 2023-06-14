import {MoreOutlined} from '@ant-design/icons';
import {Button, Space, Typography} from 'antd';

import emptyStateIcon from 'assets/SpanAssertionsEmptyState.svg';
import styled from 'styled-components';
import {ResourceType} from 'types/Resource.type';

export const ActionButton = styled(MoreOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const Container = styled.div<{$type: ResourceType}>`
  background: ${({theme}) => theme.color.white};
  border-left: ${({$type}) => ($type === ResourceType.Test ? '4px solid #2f1e61' : '4px solid #8B2C53')};
  border-radius: 2px;
  box-shadow: -1px 1px 5px #e4e9f5;
`;

export const Box = styled.div<{$type: ResourceType}>`
  align-items: center;
  background: ${({$type}) => ($type === ResourceType.Test ? '#2f1e61' : '#8B2C53')};
  border-radius: 3px;
  display: flex;
  justify-content: center;
  height: 38px;
  width: 38px;
  min-width: 38px;
  min-height: 38px;
`;

export const ResultListContainer = styled.div`
  margin: 0 68px 54px 70px;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const ButtonContainer = styled.div`
  display: flex;
  justify-content: center;
  height: 38px;
  width: 38px;
  min-width: 38px;
  min-height: 38px;
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

export const Row = styled.div<{$gap?: number; $noWrap?: boolean}>`
  align-items: center;
  column-gap: ${({$gap}) => $gap && `${$gap}px`};
  display: flex;
  white-space: ${({$noWrap}) => $noWrap && 'nowrap'};
`;

export const RunButton = styled(Button)`
  margin-right: 12px;
`;

export const RunsContainer = styled.div`
  padding: 0 24px 15px 64px;
`;

export const RunsListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

export const TestContainer = styled.div`
  cursor: pointer;
  display: grid;
  grid-template-columns: auto auto 1fr 28px 100px auto auto;
  align-items: center;
  gap: 18px;
  padding: 15px 24px;
`;

export const Text = styled(Typography.Text).attrs({as: 'p'})`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

export const TitleContainer = styled.div`
  max-width: 500px;
  min-width: 500px;
  width: 500px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin: 0;
  }
`;
