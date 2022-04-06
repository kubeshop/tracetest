import {Tabs} from 'antd';
import styled from 'styled-components';

export const TraceTabs = styled(Tabs)`
  .ant-tabs-nav {
    margin-bottom: 0;
  }

  .ant-tabs-nav::before {
    border: none;
  }
`;
