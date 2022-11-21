import styled from 'styled-components';

export const Container = styled.div<{$isTrace: boolean}>`
  padding: 24px;
  max-width: ${({$isTrace}) => ($isTrace ? 'unset' : '60vw')};
`;
