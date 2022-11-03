import {Typography} from 'antd';
import emptyStateIcon from 'assets/SpanAssertionsEmptyState.svg';
import styled from 'styled-components';

export const TestDetailsHeader = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-between;
  margin: 32px 0 24px;
`;

export const Wrapper = styled.div<{detail?: boolean}>`
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  background: ${({theme}) => theme.color.white};
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
  margin: 100px 0;
`;

export const Container = styled.div`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const Text = styled(Typography.Title).attrs({level: 3})<{opacity?: number}>`
  && {
    font-size: 12px;
    color: ${({opacity}) => `rgba(3, 24, 73, ${opacity || 1})`};
    margin: 0 !important;
  }
`;

export const TagContainer = styled.div`
  display: flex;
  white-space: nowrap;
  overflow: auto;
`;

export const Title = styled(Typography.Title).attrs({level: 3})`
  && {
    margin: 0;
  }
`;

export const Stack = styled.div`
  display: flex;
  justify-content: space-between;
  flex-direction: column;
`;

export const Info = styled.div`
  flex: 1;
`;

export const Section = styled.div`
  flex: 1;
`;

export const SectionLeft = styled.div`
  background-color: ${({theme}) => theme.color.white};
  overflow-y: auto;
  z-index: 1;
  flex-basis: 50%;
`;

export const SectionRight = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border-left: 1px solid rgba(3, 24, 73, 0.1);
  overflow-y: auto;
  z-index: 2;
  padding: 24px;
  flex-basis: 50%;
`;
