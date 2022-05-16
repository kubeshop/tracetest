import styled from 'styled-components';
import noResultsIcon from '../../assets/SpanAssertionsEmptyState.svg';

export const AssertionCardList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

export const EmptyStateContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 40px;
  flex-direction: column;
  gap: 14px;
`;

export const EmptyStateIcon = styled.img.attrs({
  src: noResultsIcon,
})``;
