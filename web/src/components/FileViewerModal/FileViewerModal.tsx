import {useCallback} from 'react';
import {Button, Typography} from 'antd';
import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import {DownloadOutlined} from '@ant-design/icons';

import * as S from './FileViewerModal.styled';
import {downloadFile} from '../../utils/Common';

interface IProps {
  data: string;
  isOpen: boolean;
  onClose(): void;
  title: string;
  language: string;
  subtitle: string;
  fileName: string;
}

const FileViewerModal = ({data, isOpen, onClose, title, language = 'javascript', subtitle, fileName}: IProps) => {
  const onDownload = useCallback(() => {
    downloadFile(data, fileName);
  }, [data, fileName]);

  const footer = (
    <>
      <Button ghost onClick={onClose} type="primary" data-cy="file-viewer-close">
        Cancel
      </Button>
      <Button data-cy="file-viewer-download" icon={<DownloadOutlined />} onClick={onDownload} type="primary">
        Download File
      </Button>
    </>
  );

  return (
    <S.FileViewerModal
      footer={footer}
      onCancel={onClose}
      title={
        <Typography.Title level={5} style={{marginBottom: 0}}>
          {title}
        </Typography.Title>
      }
      width="70%"
      visible={isOpen}
    >
      <S.SubtitleContainer>
        <Typography.Text>{subtitle}</Typography.Text>
      </S.SubtitleContainer>
      <S.CodeContainer data-cy="file-viewer-code-container">
        <SyntaxHighlighter language={language} style={arduinoLight}>
          {data}
        </SyntaxHighlighter>
      </S.CodeContainer>
    </S.FileViewerModal>
  );
};

export default FileViewerModal;
