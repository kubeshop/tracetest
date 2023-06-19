import CodeBlock from 'components/CodeBlock/CodeBlock';
import {ReadOutlined} from '@ant-design/icons';
import {CLI_RUNNING_TESTS_URL} from 'constants/Common.constants';
import * as S from './CliCommand.styled';
import Controls from './Controls';
import useCliCommand from './hooks/useCliCommand';

const CLiCommand = () => {
  const {command, onGetCommand} = useCliCommand();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>CLI Configuration</S.Title>
        <a href={CLI_RUNNING_TESTS_URL} target="_blank">
          <ReadOutlined />
        </a>
      </S.TitleContainer>
      <S.Subtitle>Tracetest CLI command:</S.Subtitle>
      <CodeBlock language="bash" value={command} minHeight="80px" maxHeight="80px" />
      <Controls onChange={onGetCommand} />
    </S.Container>
  );
};

export default CLiCommand;
