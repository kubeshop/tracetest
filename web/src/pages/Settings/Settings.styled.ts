import styled from 'styled-components';
import {Typography} from 'antd';

export const Container = styled.div`
  padding: 0 24px;
`;

export const Header = styled.div`
  padding: 23px 0;
  width: 100%;
`;

export const TabsContainer = styled.div`
  .ant-tabs-nav {
    padding: 0 12px;
    margin-bottom: 0;
  }

  .ant-tabs-nav {
    padding: 0;
  }

  .ant-tabs-content {
    margin-top: 24px;
  }
`;

export const Title = styled(Typography.Title)`
  && {
    margin: 0;
  }
`;

export const TabTextContainer = styled.div`
  display: flex;
  align-items: center;
`;
