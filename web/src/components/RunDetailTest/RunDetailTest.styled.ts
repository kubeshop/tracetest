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

export const SectionLeft = styled(Section)`
  background-color: ${({theme}) => theme.color.background};
  z-index: 1;
`;

export const SectionRight = styled(Section)`
  background-color: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
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
`;

export const CountBadge = styled(Badge)`
  .ant-badge-count {
    background-color: ${({theme}) => theme.color.backgroundDark};
    color: ${({theme}) => theme.color.text};
  }
`;

export const SpanDetailContainer = styled.div<{$isOpen: boolean}>`
  background-color: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  height: 100%;
  overflow: visible;
  position: relative;

  > div {
    opacity: 0;
    pointer-events: none;
  }

  ${({$isOpen}) =>
    $isOpen &&
    css`
      > div {
        opacity: 1;
        pointer-events: auto;
      }
    `}
`;
