import {InfoCircleFilled} from '@ant-design/icons';
import styled from 'styled-components';

export const Warning = styled(InfoCircleFilled)`
  && {
    color: ${({theme}) => theme.color.warningYellow};
  }
`;
