import {Typography} from 'antd';
import styled from 'styled-components';

import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const EmptyContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  height: calc(100% - 70px);
  justify-content: center;
  margin-top: -70px;
`;

export const EmptyIcon = styled.img.attrs({
  src: noResultsIcon,
})`
  height: auto;
  margin-bottom: 16px;
  width: 90px;
`;

export const EmptyText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;

export const EmptyTitle = styled(Typography.Title).attrs({level: 3})``;

export const SnippetsContainer = styled.div`
  margin: 16px 0;
`;
