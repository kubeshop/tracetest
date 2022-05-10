import {Tabs} from 'antd';
import styled from 'styled-components';

export const TraceTabs = styled(Tabs)`
  .ant-tabs-nav {
    margin-bottom: 0;
  }

  .ant-tabs-nav::before {
    border: none;
  }

  .ant-tabs-tabpane {
    flex-direction: column;
  }

  .ant-tabs-content-holder {
    width: 100%;
    min-width: unset;
    min-height: unset;
    overflow-y: scroll;
  }
`;

export const FailedTrace = styled.div`
  height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;
