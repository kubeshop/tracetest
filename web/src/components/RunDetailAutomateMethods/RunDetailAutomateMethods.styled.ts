import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
`;

export const TitleContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 24px;
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin: 0;
  }
`;

export const Subtitle = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textLight};
  }
`;

export const TabsContainer = styled.div`
  .ant-tabs-nav {
    padding: 0 12px;
    margin-bottom: 0;
  }

  .ant-tabs-content-holder {
    height: calc(100% - 38px);
    overflow-y: scroll;
  }

  .ant-tabs-nav {
    padding: 0;
  }
`;
