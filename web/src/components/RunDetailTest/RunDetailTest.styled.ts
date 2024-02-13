import {Badge} from 'antd';
import styled, {css} from 'styled-components';

export const Container = styled.div`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const Section = styled.div`
  flex: 1;
`;

export const SectionLeft = styled(Section)<{$isTimeline: boolean}>`
  background-color: ${({theme}) => theme.color.background};
  box-shadow: inset 20px 0px 24px -20px rgba(153, 155, 168, 0.18), inset -20px 0 24px -20px rgba(153, 155, 168, 0.18);
  z-index: 1;

  ${({$isTimeline}) =>
    $isTimeline &&
    css`
      max-size: calc(100% - 695px);
      overflow: scroll;
    `}
`;

export const SectionRight = styled(Section)`
  background-color: ${({theme}) => theme.color.white};
  overflow: hidden;
  position: relative;
  z-index: 2;
`;

export const SwitchContainer = styled.div`
  bottom: 163px;
  left: 16px;
  position: absolute;
  z-index: 9;
`;

export const TabsContainer = styled.div`
  height: 100%;
  overflow-y: auto;
  padding: 24px;
  position: relative;

  .ant-tabs-small > .ant-tabs-nav .ant-tabs-tab {
    padding: 0 0 8px;
  }

  .ant-tabs,
  .ant-tabs .ant-tabs-content {
    height: 100%;
  }
`;

export const CountBadge = styled(Badge)`
  .ant-badge-count {
    background-color: ${({theme}) => theme.color.backgroundDark};
    color: ${({theme}) => theme.color.text};
  }
`;
