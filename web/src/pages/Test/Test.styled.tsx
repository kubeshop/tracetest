import styled, {keyframes, css} from 'styled-components';
import {Tabs as AntTabs} from 'antd';

export const Header = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
  height: 56px;
  padding: 0 32px;
  border-bottom: 1px solid rgb(213, 215, 224);
`;

export const TestDetailsHeader = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-between;
  margin: 32px 0px 24px;
`;

export const Wrapper = styled.div`
  padding: 0px 24px;
`;

const IndeterminateAnimation = keyframes`
  0% {
    transform:  translateX(0) scaleX(0);
  }
  40% {
    transform:  translateX(0) scaleX(0.4);
  }
  100% {
    transform:  translateX(100%) scaleX(0.5);
  }
`;

export const TestTabs = styled(AntTabs)<{loading: string}>`
  ${props => {
    if (props.loading === 'true') {
      return css`
        &:first-child > .ant-tabs-nav {
        }
        &:first-child > .ant-tabs-nav::before,
        &:first-child > .ant-tabs-nav::after {
          z-index: 1000;
          content: '';
          height: 2px;
          bottom: 0;
          left: 0;
          position: absolute;
          width: 100%;
        }

        &:first-child > .ant-tabs-nav::before {
          background-color: #c4c4c4;
        }
        &:first-child > .ant-tabs-nav::after {
          background-color: #1890ff;
          animation: ${IndeterminateAnimation} 1s infinite linear;
          transform-origin: 0% 50%;
        }
      `;
    }
  }}
`;
