import {CloseCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const SearchContainer = styled.div`
  padding: 24px 24px 0;
  position: relative;
  z-index: 9;
`;

export const Section = styled.div`
  flex: 1;
  height: 100%;
  overflow: hidden;
  width: 100%;
  z-index: 1;
`;

export const SectionLeft = styled(Section)`
  background-color: ${({theme}) => theme.color.background};
  z-index: 1;
`;

export const SectionRight = styled(Section)<{$shouldScroll: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  overflow-y: ${({$shouldScroll}) => ($shouldScroll ? 'scroll' : 'hidden')};
  z-index: 2;
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

export const NoMatchesContainer = styled.div`
  color: ${({theme}) => theme.color.textSecondary};
  margin-top: 8px;
  margin-left: 8px;
`;

export const NoMatchesText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;
