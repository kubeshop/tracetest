import {Modal} from 'antd';
import styled from 'styled-components';

export const CodeContainer = styled.div`
  border: ${({theme}) => `1px solid ${theme.color.border}`};
  min-height: 370px;

  pre {
    margin: 0;
    min-height: inherit;
  }
`;

export const FileViewerModal = styled(Modal)`
  & .ant-modal-body {
    background: ${({theme}) => theme.color.bg};
  }
`;

export const SubtitleContainer = styled.div`
  margin-bottom: 8px;
`;
