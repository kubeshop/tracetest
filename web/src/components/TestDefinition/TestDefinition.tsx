import {useCallback, useEffect} from 'react';
import {snakeCase} from 'lodash';
import {DownloadOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import {downloadFile} from 'utils/Common';
import useCopy from 'hooks/useCopy';
import useDefinitionFile from 'hooks/useDefinitionFile';
import {ResourceType} from 'types/Resource.type';
import Test from 'models/Test.model';
import * as S from './TestDefinition.styled';

interface IProps {
  test: Test;
}

const TestDefinition = ({test: {id, version, name}}: IProps) => {
  const {definition, loadDefinition} = useDefinitionFile();

  const onDownload = useCallback(() => {
    downloadFile(definition, `${snakeCase(`${name} ${version} definition`)}.yaml`);
  }, [definition, name, version]);

  const copy = useCopy();

  useEffect(() => {
    loadDefinition(ResourceType.Test, id, version);
  }, [id, loadDefinition, version]);

  return (
    <S.Container>
      <S.Title>Test Definition</S.Title>
      <S.SubtitleContainer>
        <Typography.Text>Preview your YAML file</Typography.Text>
      </S.SubtitleContainer>
      <S.CodeContainer data-cy="test-definition-code-container">
        <S.CopyIconContainer>
          <S.CopyButton onClick={() => copy(definition)} ghost type="primary">
            Copy
          </S.CopyButton>
        </S.CopyIconContainer>
        <SyntaxHighlighter language="yaml" style={arduinoLight}>
          {definition}
        </SyntaxHighlighter>
      </S.CodeContainer>
      <S.Footer>
        <Button data-cy="file-viewer-download" icon={<DownloadOutlined />} onClick={onDownload} type="primary">
          Download File
        </Button>
      </S.Footer>
    </S.Container>
  );
};

export default TestDefinition;
