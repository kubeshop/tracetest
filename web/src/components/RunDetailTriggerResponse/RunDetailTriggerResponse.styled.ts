import {Typography} from 'antd';
import styled from 'styled-components';

export {default as AttributeTitle} from 'components/AttributeRow/AttributeTitle';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 24px;
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
    font-weight: 700;
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
    margin-left: 14px;
    font-size: ${({theme}) => theme.size.md};
  }
`;

export const LoadingResponseBody = styled.div`
  margin-top: 25px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  height: 100px;
  padding: 0.4em 0.6em;
  background: ${({theme}) => theme.color.background};
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  font-size: ${({theme}) => theme.size.sm};
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
  height: 8px;
  border-radius: 2px;
  width: ${({$width = 100}) => $width}%;
`;

export const ValueJson = styled(Typography.Text)`
  margin-top: 25px;
  display: block;

  pre {
    margin: 0;
    background: ${({theme}) => theme.color.background};
    border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    font-size: ${({theme}) => theme.size.sm};
  }
`;

export const HeaderContainer = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  margin-bottom: 4px;
  padding: 12px;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: ${({theme}) => theme.color.background};
  }
`;

export const Header = styled.div`
  flex: 1;
`;

export const AttributeValueRow = styled.div`
  display: flex;
  word-break: break-word;
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
  font-weight: 700;
`;
