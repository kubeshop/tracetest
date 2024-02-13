import {CloseCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div<{$isTimeline: boolean}>`
  display: flex;
  height: 100%;
  width: 100%;
  min-width: ${({$isTimeline}) => $isTimeline && '1000px'};
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

export const SectionLeft = styled(Section)<{$hasShadow?: boolean}>`
  background-color: ${({theme}) => theme.color.background};
  box-shadow: ${({$hasShadow}) =>
    $hasShadow &&
    `inset 20px 0px 24px -20px rgba(153, 155, 168, 0.18), inset -20px 0 24px -20px rgba(153, 155, 168, 0.18)`};
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

export const NoMatchesContainer = styled.div`
  color: ${({theme}) => theme.color.textSecondary};
  margin-top: 8px;
  margin-left: 8px;
`;

export const NoMatchesText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;
