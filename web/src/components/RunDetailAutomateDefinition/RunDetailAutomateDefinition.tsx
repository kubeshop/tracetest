import {DownloadOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {useCallback} from 'react';
import {FramedCodeBlock} from 'components/CodeBlock';
import {Overlay} from 'components/Inputs';
import {ResourceName, ResourceType} from 'types/Resource.type';
import {downloadFile} from 'utils/Common';
import * as S from './RunDetailAutomateDefinition.styled';

interface IProps {
  resourceType: ResourceType;
  fileName: string;
  definition: string;
  onFileNameChange(value: string): void;
}

const RunDetailAutomateDefinition = ({resourceType, fileName, definition, onFileNameChange}: IProps) => {
  const onDownload = useCallback(() => {
    downloadFile(definition, fileName);
  }, [definition, fileName]);

  return (
    <S.Container>
      <S.Title>{ResourceName[resourceType]} Definition</S.Title>
      <S.FileName>
        <Overlay value={fileName} onChange={onFileNameChange} />
      </S.FileName>

      <FramedCodeBlock
        title="Preview your YAML file"
        value={definition}
        language="yaml"
        actions={
          <Button
            data-cy="file-viewer-download"
            icon={<DownloadOutlined />}
            onClick={onDownload}
            size="small"
            type="primary"
          >
            Download File
          </Button>
        }
        isFullHeight
      />
    </S.Container>
  );
};

export default RunDetailAutomateDefinition;
