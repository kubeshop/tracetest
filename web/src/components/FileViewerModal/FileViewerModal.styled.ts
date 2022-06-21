import {Modal} from 'antd';
import styled from 'styled-components';

export const CodeContainer = styled.div`
  border: 1px solid rgba(3, 24, 73, 0.1);
  min-height: 370px;

  pre {
    margin: 0;
    min-height: inherit;
  }
`;

export const FileViewerModal = styled(Modal)`
  & .ant-modal-body {
    background: #fbfbff;
  }
`;

export const SubtitleContainer = styled.div`
  margin-bottom: 8px;
`;
