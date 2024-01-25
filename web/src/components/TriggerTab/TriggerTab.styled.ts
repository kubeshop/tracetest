import {Badge as AntdBadge, Typography} from 'antd';
import styled from 'styled-components';

export const Badge = styled(AntdBadge)`
  .ant-badge-status-text {
    margin-left: 0;
  }
`;

export const Text = styled(Typography.Text)`
  color: ${({theme}) => theme.color.success};
`;
