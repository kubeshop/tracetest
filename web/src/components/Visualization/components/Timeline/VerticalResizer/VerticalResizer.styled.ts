import styled from 'styled-components';

export const VerticalResizer = styled.div`
  left: 0;
  position: absolute;
  right: 0;
  top: 0;
`;

export const VerticalResizerDragger = styled.div`
  border-left: 1px solid rgb(222, 227, 236);
  cursor: ew-resize;
  height: calc(100vh - 50px);
  margin-left: -1px;
  position: absolute;
  top: 0;
  width: 1px;
  z-index: 2;

  ::before {
    position: absolute;
    top: 0;
    bottom: 0;
    left: -8px;
    right: 0;
    content: ' ';
  }

  :hover {
    border-left: 2px solid ${({theme}) => theme.color.border};
  }

  &.dragging {
    background: rgba(136, 0, 136, 0.05);
    width: unset;

    ::before {
      left: -2000px;
      right: -2000px;
    }
  }
`;
