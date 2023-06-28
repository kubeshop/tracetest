import {FramedCodeBlock} from 'components/CodeBlock';
import {ReadOutlined} from '@ant-design/icons';
import {CLI_RUNNING_TESTS_URL} from 'constants/Common.constants';
import * as S from './CliCommand.styled';
import Controls from './Controls';
import useCliCommand from './hooks/useCliCommand';
import {IMethodProps} from '../../RunDetailAutomateMethods';

const CLiCommand = ({test, environmentId, fileName = ''}: IMethodProps) => {
  const {command, onGetCommand} = useCliCommand();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>CLI Configuration</S.Title>
        <a href={CLI_RUNNING_TESTS_URL} target="_blank">
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
      <Controls onChange={onGetCommand} test={test} fileName={fileName} environmentId={environmentId} />
    </S.Container>
  );
};

export default CLiCommand;
