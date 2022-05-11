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

export const Main = styled.main<{height: string}>`
  display: flex;
  width: 100%;
  min-height: ${({height}) => height};
  max-height: ${({height}) => height};
  height: ${({height}) => height};
`;

export const DiagramSection = styled.div`
  flex-basis: 50%;
  display: flex;
  flex-direction: column;
  padding: 24px;
`;

export const DetailsSection = styled.div`
  flex-basis: 50%;
  overflow-y: scroll;
  background: #fff;
  box-shadow: 0px 20px 24px rgba(153, 155, 168, 0.18);
`;

export const TabsContainer = styled.div`
  padding: 14px 24px;
  overflow: hidden;
`;
