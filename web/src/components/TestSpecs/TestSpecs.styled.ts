import {Typography} from 'antd';
import styled from 'styled-components';

import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

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

export const ListContainer = styled.div`
  flex: 1;
`;

export const HiddenElementContainer = styled.div`
  position: absolute;
  visibility: hidden;
  z-index: -1;
`;

export const HiddenElement = styled.div`
  font-family: monospace;
  font-size: ${({theme}) => theme.size.md};
  line-height: 1.4;
  padding: 4px 96px 76px 4px;
  visibility: hidden;
  width: 100%;
`;
