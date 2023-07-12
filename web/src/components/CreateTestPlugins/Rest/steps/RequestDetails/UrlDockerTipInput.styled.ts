import {Typography} from 'antd';
import styled from 'styled-components';

export const Paragraph = styled(Typography.Paragraph)`
  && {
    margin-top: 10px;
    font-weight: 600;

    span {
      font-weight: 400;
    }
  }
`;
