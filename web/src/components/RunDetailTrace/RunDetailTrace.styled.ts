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
  overflow: hidden;
  width: 100%;
  z-index: 1;
`;

export const VisualizationContainer = styled.div`
  height: calc(100% - 80px);
  position: relative;
`;

export const SwitchContainer = styled.div<{$hasSpace: boolean}>`
  position: absolute;
  z-index: 9;
  left: 16px;
  bottom: 163px;
`;
