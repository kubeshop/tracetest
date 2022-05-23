import styled, {css} from 'styled-components';

export const Main = styled.main<{height: string}>`
  display: flex;
  width: 100%;
  min-height: ${({height}) => height};
  max-height: ${({height}) => height};
  height: ${({height}) => height};
`;

export const DiagramSection = styled.div`
  flex-basis: 50%;
  display: flex;
  flex-direction: column;
  padding: 24px;
`;

export const DetailsSection = styled.div`
  flex-basis: 50%;
  overflow-y: scroll;
  background: #fff;
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
`;

export const Header = styled.div<{visiblePortion: number}>`
  display: flex;
  align-items: center;
  cursor: pointer;
  justify-content: space-between;
  width: 100%;
  background: #f5f5fa;
  ${props =>
    css`
      height: ${props.visiblePortion}px;
    `}
  padding: 0 24px;
  color: rgb(213, 215, 224);
`;
