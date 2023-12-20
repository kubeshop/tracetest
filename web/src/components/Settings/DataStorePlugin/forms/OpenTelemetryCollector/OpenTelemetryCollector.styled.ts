import {CopyOutlined} from '@ant-design/icons';
import {Modal, Typography} from 'antd';
import styled from 'styled-components';

export const CollapseContainer = styled.div`
  .ant-collapse > .ant-collapse-item > .ant-collapse-header {
    padding: 12px 16px 12px 0;
  }

  .ant-collapse-content > .ant-collapse-content-box {
    padding: 12px 16px 12px 0;
  }
`;

export const CodeContainer = styled.div`
  margin-bottom: 18px;
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

export const Container = styled.div``;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.md};
  }
`;

export const Description = styled(Typography.Paragraph)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
    font-size: ${({theme}) => theme.size.md};
  }
`;

export const SwitchContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 18px;
`;

export const UrlEntry = styled.div`
  display: grid;
  grid-template-columns: auto 1fr;
  align-items: center;
  gap: 16px;
`;
