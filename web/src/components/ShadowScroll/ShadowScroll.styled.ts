import styled from 'styled-components';

export const Overlay = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  background: transparent;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
  opacity: 0;
`;

export const Shadow = styled(Overlay)<{$direction: 'to left' | 'to right' | 'to top' | 'to bottom'}>`
  background: linear-gradient(${({$direction}) => $direction}, #ffffff 6.83%, rgba(255, 255, 255, 0) 25%);
`;

export const Content = styled.div`
  min-width: 100%;
  max-width: 100%;
  overflow-x: scroll;

  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
    -webkit-appearance: none;
    width: 0;
    height: 0;
  }
`;
