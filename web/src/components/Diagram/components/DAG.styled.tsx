import styled, {css} from 'styled-components';

export const Container = styled.div<{showAffected: boolean}>`
  position: relative;
  height: 100%;

  ${({showAffected}) =>
    showAffected &&
    css`
      .react-flow__node-TraceNode:not(.affected) > div {
        opacity: 0.5;
      }
    `}
`;
