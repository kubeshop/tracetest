import styled, {css} from 'styled-components';
import noResultsIcon from '../../../assets/SpanAssertionsEmptyState.svg';

export const Header = styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
`;

export const TraceHeader = styled.div<{visiblePortion: number}>`
  display: flex;
  align-items: center;
  align-content: center;
  align-items: center;
  cursor: grab;
  width: 100%;
  ${props =>
    css`
      height: ${props.visiblePortion}px;
    `}
  margin: 0 16px;
  color: rgb(213, 215, 224);
`;

export const Container = styled.div`
  padding: 24px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  margin-bottom: 16px;
  min-height: 280px;
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
