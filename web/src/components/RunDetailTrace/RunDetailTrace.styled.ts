import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const SearchContainer = styled.div`
  padding: 24px 24px 0;
`;

export const Section = styled.div`
  flex: 1;
  z-index: 1;
  width: 100%;
`;

export const VisualizationContainer = styled.div`
  height: calc(100% - 80px);
  position: relative;
`;

export const SwitchContainer = styled.div<{$hasSpace: boolean}>`
  position: absolute;
  right: ${({$hasSpace}) => ($hasSpace ? '24px' : '132px')};
  top: 16px;
  z-index: 9;
`;
