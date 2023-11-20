import styled from 'styled-components';
import {Modal as AntModal, Typography} from 'antd';

export const Modal = styled(AntModal)`
  min-width: 625px;

  .ant-modal-body {
    background: ${({theme}) => theme.color.background};
    overflow-y: scroll;
    max-height: calc(100vh - 250px);
    position: relative;
  }
`;

export const Container = styled.div`
  min-height: 361px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 0px;
  }
`;
