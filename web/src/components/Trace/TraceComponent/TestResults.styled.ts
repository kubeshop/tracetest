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
  overflow-y: auto;
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

export const TraceContainer = styled.div<{height: string}>`
  display: flex;
  width: 100%;
  min-height: ${props => props.height};
  max-height: ${props => props.height};
  height: ${props => props.height};
`;

export const LeftContainer = styled.div`
  flex-basis: 50%;
  padding-top: 10px;
  padding-left: 10px;
`;

export const RightContainer = styled.div`
  display: flex;
  flex-basis: 50%;
  width: 100%;
  flex-grow: 1;
  padding-top: 10px;
  padding-left: 10px;
`;
