import styled from 'styled-components';
import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const EmptyStateContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 14px;
  justify-content: center;
  margin-top: 40px;
`;

export const EmptyStateIcon = styled.img.attrs({
  src: noResultsIcon,
})`
  height: auto;
  width: 90px;
`;
