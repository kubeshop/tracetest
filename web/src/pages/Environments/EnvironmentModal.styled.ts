import {Modal as AntModal, Typography} from 'antd';
import styled from 'styled-components';

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

export const Modal = styled(AntModal)`
  .ant-modal-body {
    background: ${({theme}) => theme.color.background};
  }
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin: 0;
  }
`;
