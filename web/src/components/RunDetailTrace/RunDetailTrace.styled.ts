import {CloseCircleFilled} from '@ant-design/icons';
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
  height: 100%;
  overflow: hidden;
  width: 100%;
  z-index: 1;
`;

export const VisualizationContainer = styled.div`
  height: calc(100% - 52px);
  position: relative;
`;

export const SwitchContainer = styled.div`
  bottom: 163px;
  left: 16px;
  position: absolute;
  z-index: 9;
`;

export const ClearSearchIcon = styled(CloseCircleFilled)`
  position: absolute;
  right: 8px;
  top: 8px;
  color: ${({theme}) => theme.color.textLight};
  cursor: pointer;
`;
