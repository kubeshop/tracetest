import {Typography} from 'antd';
import styled from 'styled-components';

export const Paragraph = styled(Typography.Paragraph)`
  && {
    margin: 0;
    font-weight: 600;

    span {
      font-weight: 400;
    }
  }
`;
