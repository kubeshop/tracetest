import {Typography} from 'antd';
import styled from 'styled-components';
import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
`;

export const TitleContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin: 0;
  }
`;

export const TabsContainer = styled.div`
  .ant-tabs-nav {
    padding: 0 12px;
    margin-bottom: 0;
  }

  .ant-tabs-content-holder {
    height: calc(100% - 38px);
    overflow-y: scroll;
  }

  .ant-tabs-nav {
    padding: 0;
  }
`;

export const StatusText = styled(Typography.Text)`
  && {
    font-size: ${({theme}) => theme.size.md};
  }
`;

export const LoadingResponseBody = styled.div`
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  justify-items: center;
  gap: 16px;
  padding: 0.4em 0.6em;
`;

export const TextHolder = styled.div<{$width?: number}>`
  @keyframes skeleton-loading {
    0% {
      background-color: hsl(200, 20%, 80%);
    }
    100% {
      background-color: hsl(200, 20%, 95%);
    }
  }

  animation: skeleton-loading 1s linear infinite alternate;
  height: 16px;
  border-radius: 2px;
  width: ${({$width = 100}) => $width}%;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
`;

export const StatusSpan = styled.span<{$isError: boolean}>`
  color: ${({$isError, theme}) => ($isError ? theme.color.error : theme.color.success)};
  font-weight: 600;
`;

export const HeadersList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px 0;
`;

export const Actions = styled.div`
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-top: 16px;
  gap: 10px;
`;

export const ResponseVarsContainer = styled.div`
  padding: 16px 0;
`;

export const EmptyContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-top: 100px;
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

export const ResponseBodyContainer = styled.div`
  display: grid;
  grid-template-columns: 95% 5%;
  width: 100%;
`;

export const ResponseBodyContent = styled.div`
  flex: 1;
  margin-top: 16px;
`;

export const ResponseBodyActions = styled.div`
  margin: 16px 0 0 4px;
`;
