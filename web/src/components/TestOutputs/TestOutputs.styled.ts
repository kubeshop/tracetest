import {Typography} from 'antd';
import styled from 'styled-components';

import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

export const Container = styled.div`
  height: 100%;
  position: relative;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  justify-content: flex-end;
  margin-bottom: 38px;
`;

export const ListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const Actions = styled.div`
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-top: 16px;
  gap: 10px;
`;

export const EmptyContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  height: calc(100% - 70px);
  justify-content: center;
  margin-top: -70px;
  text-align: center;
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
