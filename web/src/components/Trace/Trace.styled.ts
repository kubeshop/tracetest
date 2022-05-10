import {Tabs} from 'antd';
import styled from 'styled-components';

export const TraceTabs = styled(Tabs)`
  flex-grow: 1;

  .ant-tabs {
    flex-grow: 1;
  }

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
    overflow-y: auto;
  }
`;
