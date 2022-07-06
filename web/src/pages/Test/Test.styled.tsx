import styled from 'styled-components';

import emptyStateIcon from 'assets/SpanAssertionsEmptyState.svg';

export const TestDetailsHeader = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-between;
  margin: 32px 0px 24px;
`;

export const Wrapper = styled.div<{detail?: boolean}>`
  padding: 0px 24px;
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
