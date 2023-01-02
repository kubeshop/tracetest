import {CopyOutlined} from '@ant-design/icons';
import {Modal} from 'antd';
import styled from 'styled-components';

export const CodeContainer = styled.div`
  position: relative;
  border: ${({theme}) => `1px solid ${theme.color.border}`};
  max-height: 370px;
  overflow-y: scroll;

  pre {
    margin: 0;
    min-height: inherit;
  }
`;

export const FileViewerModal = styled(Modal)`
  & .ant-modal-body {
    background: ${({theme}) => theme.color.background};
  }
`;

export const CopyIcon = styled(CopyOutlined)`
  color: ${({theme}) => theme.color.primary};
`;

export const CopyIconContainer = styled.div`
  position: absolute;
  right: 8px;
  top: 9px;
  padding: 0 2px;
  border-radius: 2px;
  cursor: pointer;
  background-color: ${({theme}) => theme.color.textHighlight};
  z-index: 101;
`;

export const SubtitleContainer = styled.div`
  margin-bottom: 8px;
`;
