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
`;

export const FailedTrace = styled.div`
  height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;
