import {DownloadOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {useCallback, useEffect} from 'react';
import {FramedCodeBlock} from 'components/CodeBlock';
import InputOverlay from 'components/InputOverlay/InputOverlay';
import useDefinitionFile from 'hooks/useDefinitionFile';
import {ResourceName, ResourceType} from 'types/Resource.type';
import {downloadFile} from 'utils/Common';
import * as S from './RunDetailAutomateDefinition.styled';

interface IProps {
  id: string;
  version: number;
  resourceType: ResourceType;
  fileName: string;
  onFileNameChange(value: string): void;
}

const RunDetailAutomateDefinition = ({id, version, resourceType, fileName, onFileNameChange}: IProps) => {
  const {definition, loadDefinition} = useDefinitionFile();

  const onDownload = useCallback(() => {
    downloadFile(definition, fileName);
  }, [definition, fileName]);

  useEffect(() => {
    loadDefinition(resourceType, id, version);
  }, [id, loadDefinition, resourceType, version]);

  return (
    <S.Container>
      <S.Title>{ResourceName[resourceType]} Definition</S.Title>
      <S.FileName>
        <InputOverlay value={fileName} onChange={onFileNameChange} />
      </S.FileName>
      <FramedCodeBlock title="Preview your YAML file" value={definition} language="yaml" />
      <S.Footer>
        <Button data-cy="file-viewer-download" icon={<DownloadOutlined />} onClick={onDownload} type="primary">
          Download File
        </Button>
      </S.Footer>
    </S.Container>
  );
};

export default RunDetailAutomateDefinition;
