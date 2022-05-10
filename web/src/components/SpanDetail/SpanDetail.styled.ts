import styled from 'styled-components';
import noResultsIcon from '../../assets/SpanAssertionsEmptyState.svg';

export const DetailsHeader = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
`;

export const DetailsContainer = styled.div`
  padding: 24px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  margin-bottom: 16px;
`;

export const DetailsEmptyStateContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 40px;
  flex-direction: column;
  gap: 14px;
  overflow-y: auto;
`;

export const DetailsTableEmptyStateIcon = styled.img.attrs({
  src: noResultsIcon,
})``;
