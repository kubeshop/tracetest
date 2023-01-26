import {Typography} from 'antd';
import styled from 'styled-components';

export const Mark = styled(Typography.Text)`
  && {
    cursor: pointer;
    font-size: ${({theme}) => theme.size.xs};
  }
`;
