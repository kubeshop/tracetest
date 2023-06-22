import {Modal} from 'antd';
import styled from 'styled-components';

export const FileViewerModal = styled(Modal)`
  top: 50px;

  & .ant-modal-body {
    background: ${({theme}) => theme.color.background};
    max-height: calc(100vh - 250px);
    overflow: scroll;
  }
`;
