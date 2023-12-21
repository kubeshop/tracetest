import {Button, Typography} from 'antd';
import styled, {createGlobalStyle} from 'styled-components';

export const StopContainer = styled.div`
  margin-left: 12px;
`;

export const GlobalStyle = createGlobalStyle`
  .ant-popover.ant-popover-placement-bottomRight {
    z-index: 9999;
  }

  #skip-trace-popover {
    .ant-popover-title {
      padding: 14px;
      border: 0;
      padding-bottom: 0;
    }

    .ant-popover-inner-content {
      padding: 14px;
      padding-top: 0;
      max-width: 520px;
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

export const ForwardButton = styled(Button)`
  && {
    font-size: ${({theme}) => theme.size.xl};
  }
`;

export const ContentContainer = styled.div`
  display: flex;
  gap: 12px;
  justify-content: space-between;
  align-items: center;
`;

export const CloseIcon = styled(Typography.Text)`
  && {
    cursor: pointer;
    color: ${({theme}) => theme.color.textSecondary};
  }
`;
