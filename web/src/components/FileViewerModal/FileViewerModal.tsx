import {useCallback} from 'react';
import {Button, Typography} from 'antd';
import {DownloadOutlined} from '@ant-design/icons';
import {downloadFile} from 'utils/Common';
import * as S from './FileViewerModal.styled';
import {FramedCodeBlock} from '../CodeBlock';

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
        <Typography.Title level={2} style={{marginBottom: 0}}>
          {title}
        </Typography.Title>
      }
      width="70%"
      visible={isOpen}
    >
      <div data-cy="file-viewer-code-container">
        <FramedCodeBlock title={subtitle} value={data} language={language} />
      </div>
    </S.FileViewerModal>
  );
};

export default FileViewerModal;
