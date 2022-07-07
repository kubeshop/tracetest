import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div<{hasMinWidth?: boolean}>`
  display: flex;
  flex-direction: column;
  height: 28px;
  justify-content: center;
  max-width: 204px;
  min-width: ${({hasMinWidth}) => hasMinWidth && '115px'};

  .ant-progress {
    font-size: 6px;
    line-height: 6px;
  }
`;

export const Text = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.sm};
  margin-top: 2px;
`;
