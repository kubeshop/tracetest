import {Typography} from 'antd';
import styled, {createGlobalStyle} from 'styled-components';

export const StopContainer = styled.div`
  margin-left: 12px;
`;

export const GlobalStyle = createGlobalStyle`
  #skip-trace-popover {
    .ant-popover-title {
      padding: 14px;
      border: 0;
      padding-bottom: 0;
    }

    .ant-popover-inner-content {
      padding: 5px 14px;
      padding-top: 0;
    }
  }
`;

export const Actions = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: space-between;
  margin-top: 24px;
`;

export const Title = styled(Typography.Title).attrs({
  level: 3,
})`
  && {
    margin: 0;
  }
`;
