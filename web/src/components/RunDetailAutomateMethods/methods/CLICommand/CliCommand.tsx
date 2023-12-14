import {ReadOutlined} from '@ant-design/icons';
import {FramedCodeBlock} from 'components/CodeBlock';
import * as S from './CliCommand.styled';
import Controls from './Controls';
import useCliCommand from './hooks/useCliCommand';
import {IMethodChildrenProps} from '../../RunDetailAutomateMethods';

const CLiCommand = ({id, variableSetId, fileName = '', resourceType, docsUrl}: IMethodChildrenProps) => {
  const {command, onGetCommand} = useCliCommand();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>CLI Configuration</S.Title>
        <a href={docsUrl} target="_blank">
          <ReadOutlined />
        </a>
      </S.TitleContainer>
      <FramedCodeBlock
        title="Tracetest CLI command:"
        language="bash"
        value={command}
        minHeight="100px"
        maxHeight="100px"
      />
      <Controls
        onChange={onGetCommand}
        id={id}
        fileName={fileName}
        variableSetId={variableSetId}
        resourceType={resourceType}
      />
    </S.Container>
  );
};

export default CLiCommand;
